package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/stripe/stripe-go/v82"
	billingportal "github.com/stripe/stripe-go/v82/billingportal/session"
	checkout "github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/subscription"
	"github.com/stripe/stripe-go/v82/subscriptionitem"
	"github.com/stripe/stripe-go/v82/webhook"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/config"
	"github.com/trysourcetool/onprem-portal/internal/core"
	"github.com/trysourcetool/onprem-portal/internal/database"
	"github.com/trysourcetool/onprem-portal/internal/errdefs"
)

type createCheckoutSessionRequest struct {
	PlanID string `json:"planId" validate:"required"`
}

type createCheckoutSessionResponse struct {
	URL string `json:"url"`
}

func (s *Server) handleCreateCheckoutSession(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctxUser := internal.ContextUser(ctx)
	var req createCheckoutSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}
	if err := validateRequest(req); err != nil {
		return err
	}
	planUUID, err := uuid.FromString(req.PlanID)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}
	plan, err := s.db.Plan().GetByID(ctx, planUUID)
	if err != nil {
		return err
	}
	sub, err := s.db.Subscription().GetByUserID(ctx, ctxUser.ID)
	if err != nil {
		return err
	}
	stripe.Key = config.Config.Stripe.Key
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(plan.StripePriceID),
				Quantity: stripe.Int64(max(sub.SeatCount, 1)),
			},
		},
		AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{
			Enabled: stripe.Bool(true),
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id": ctxUser.ID.String(),
			},
		},
		SuccessURL: stripe.String(config.Config.BaseURL + "/settings/billing"),
		CancelURL:  stripe.String(config.Config.BaseURL + "/settings/billing"),
	}
	sess, err := checkout.New(params)
	if err != nil {
		return errdefs.ErrInternal(err)
	}
	return s.renderJSON(w, http.StatusOK, createCheckoutSessionResponse{URL: sess.URL})
}

type customerPortalUrlResponse struct {
	URL string `json:"url"`
}

func (s *Server) handleGetCustomerPortalUrl(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctxUser := internal.ContextUser(ctx)
	sub, err := s.db.Subscription().GetByUserID(ctx, ctxUser.ID)
	if err != nil {
		return err
	}
	if sub.StripeCustomerID == "" {
		return s.renderJSON(w, http.StatusBadRequest, statusResponse{Code: http.StatusBadRequest, Message: "No Stripe customer found"})
	}
	stripe.Key = config.Config.Stripe.Key
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(sub.StripeCustomerID),
		ReturnURL: stripe.String(config.Config.BaseURL + "/settings/billing"),
	}
	portal, err := billingportal.New(params)
	if err != nil {
		return errdefs.ErrInternal(err)
	}
	return s.renderJSON(w, http.StatusOK, customerPortalUrlResponse{URL: portal.URL})
}

func (s *Server) handleStripeWebhook(w http.ResponseWriter, r *http.Request) error {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return s.renderJSON(w, http.StatusServiceUnavailable, statusResponse{Code: http.StatusServiceUnavailable, Message: "Read error"})
	}
	sigHeader := r.Header.Get("Stripe-Signature")
	endpointSecret := config.Config.Stripe.WebhookSecret
	if endpointSecret == "" {
		return s.renderJSON(w, http.StatusInternalServerError, statusResponse{Code: http.StatusInternalServerError, Message: "Webhook secret not configured"})
	}
	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		return s.renderJSON(w, http.StatusBadRequest, statusResponse{Code: http.StatusBadRequest, Message: "Signature verification failed"})
	}

	switch event.Type {
	case "customer.subscription.created":
		var subObj stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subObj); err != nil {
			break
		}
		if subObj.Items == nil || subObj.Items.Data == nil || len(subObj.Items.Data) == 0 {
			break
		}
		ctx := r.Context()
		userID := subObj.Metadata["user_id"]
		if userID == "" {
			break
		}
		userUUID, err := uuid.FromString(userID)
		if err != nil {
			break
		}
		u, err := s.db.User().GetByID(ctx, userUUID)
		if err != nil {
			break
		}
		plan, err := s.db.Plan().GetByStripePriceID(ctx, subObj.Items.Data[0].Price.ID)
		if err != nil {
			break
		}
		sub, err := s.db.Subscription().GetByUserID(ctx, u.ID)
		if err != nil {
			break
		}
		sub.Status = core.SubscriptionStatusActive
		sub.StripeCustomerID = subObj.Customer.ID
		sub.StripeSubscriptionID = subObj.ID
		sub.PlanID = &plan.ID
		if err := s.db.Subscription().Update(ctx, sub); err != nil {
			break
		}
	case "customer.subscription.updated":
		var subObj stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subObj); err != nil {
			break
		}
		ctx := r.Context()
		plan, err := s.db.Plan().GetByStripePriceID(ctx, subObj.Items.Data[0].Price.ID)
		if err != nil {
			break
		}
		sub, err := s.db.Subscription().GetByStripeSubscriptionID(ctx, subObj.ID)
		if err != nil {
			break
		}
		sub.PlanID = &plan.ID
		// Update status based on Stripe subscription status
		switch subObj.Status {
		case "active":
			sub.Status = core.SubscriptionStatusActive
		case "trialing":
			sub.Status = core.SubscriptionStatusTrial
		case "canceled":
			sub.Status = core.SubscriptionStatusCanceled
			_ = s.scheduleUserDeletion(ctx, sub.UserID)
		case "past_due":
			sub.Status = core.SubscriptionStatusPastDue
			_ = s.scheduleUserDeletion(ctx, sub.UserID)
		default:
			sub.Status = core.SubscriptionStatusUnknown
		}
		if err := s.db.Subscription().Update(ctx, sub); err != nil {
			break
		}
	case "customer.subscription.deleted":
		var subObj stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subObj); err != nil {
			break
		}
		ctx := r.Context()
		sub, err := s.db.Subscription().GetByStripeSubscriptionID(ctx, subObj.ID)
		if err != nil {
			break
		}
		sub.Status = core.SubscriptionStatusCanceled
		if err := s.db.Subscription().Update(ctx, sub); err != nil {
			break
		}
		_ = s.scheduleUserDeletion(ctx, sub.UserID)
	default:
		// ignore other events
	}
	return s.renderJSON(w, http.StatusOK, statusResponse{Code: http.StatusOK, Message: "ok"})
}

type updateSeatsRequest struct {
	Seats int64 `json:"seats" validate:"required"`
}

func (s *Server) handleUpdateSeats(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req updateSeatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}
	if err := validateRequest(req); err != nil {
		return err
	}
	license := internal.ContextLicense(ctx)
	if license == nil {
		return errdefs.ErrUnauthenticated(fmt.Errorf("license not found in context"))
	}
	sub, err := s.db.Subscription().GetByUserID(ctx, license.UserID)
	if err != nil {
		return err
	}

	stripeSub, err := subscription.Get(sub.StripeSubscriptionID, nil)
	if err != nil {
		return errdefs.ErrInternal(err)
	}
	if stripeSub.Items == nil || len(stripeSub.Items.Data) == 0 {
		return errdefs.ErrInternal(fmt.Errorf("stripe subscription items not found for user"))
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		seatCount, err := tx.Subscription().UpdateSeatCount(ctx, sub.ID, req.Seats)
		if err != nil {
			return err
		}

		if sub.IsTrial() {
			return nil
		}
		if sub.IsPastDue() {
			return errdefs.ErrInternal(fmt.Errorf("subscription is past due"))
		}
		if sub.IsCanceled() {
			return errdefs.ErrInternal(fmt.Errorf("subscription is canceled"))
		}
		if sub.StripeCustomerID == "" {
			return errdefs.ErrInternal(fmt.Errorf("stripe customer id not found for user"))
		}
		if sub.StripeSubscriptionID == "" {
			return errdefs.ErrInternal(fmt.Errorf("stripe subscription id not found for user"))
		}

		stripe.Key = config.Config.Stripe.Key
		params := &stripe.SubscriptionItemParams{
			Quantity: stripe.Int64(seatCount),
		}
		_, err = subscriptionitem.Update(stripeSub.Items.Data[0].ID, params)
		if err != nil {
			return errdefs.ErrInternal(err)
		}

		return nil
	}); err != nil {
		return err
	}
	return s.renderJSON(w, http.StatusOK, statusResponse{Code: 0, Message: "ok"})
}

func (s *Server) scheduleUserDeletion(ctx context.Context, userID uuid.UUID) error {
	u, err := s.db.User().GetByID(ctx, userID)
	if err != nil {
		return err
	}
	u.SetScheduledDeletionAt()
	return s.db.User().Update(ctx, u)
}
