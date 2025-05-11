package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/stripe/stripe-go/v82"
	billingportal "github.com/stripe/stripe-go/v82/billingportal/session"
	checkout "github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/webhook"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/config"
	"github.com/trysourcetool/onprem-portal/internal/core"
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
	customerID := sub.StripeCustomerID
	if customerID == "" {
		cusParams := &stripe.CustomerParams{
			Email: stripe.String(ctxUser.Email),
			Name:  stripe.String(ctxUser.FullName()),
		}
		cusObj, err := customer.New(cusParams)
		if err != nil {
			return errdefs.ErrInternal(err)
		}
		customerID = cusObj.ID
		sub.StripeCustomerID = customerID
		if err := s.db.Subscription().Update(ctx, sub); err != nil {
			return errdefs.ErrInternal(err)
		}
	}
	params := &stripe.CheckoutSessionParams{
		Customer:           stripe.String(customerID),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(plan.StripePriceID),
				Quantity: stripe.Int64(1),
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
	payload, err := ioutil.ReadAll(r.Body)
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
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			break
		}
		if session.Subscription == nil || session.Customer == nil {
			break
		}
		ctx := r.Context()
		sub, err := s.db.Subscription().GetByStripeSubscriptionID(ctx, session.Subscription.ID)
		if err != nil {
			break
		}
		sub.Status = core.SubscriptionStatusActive
		sub.StripeCustomerID = session.Customer.ID
		sub.StripeSubscriptionID = session.Subscription.ID
		if err := s.db.Subscription().Update(ctx, sub); err != nil {
			break
		}
	case "customer.subscription.updated":
		var subObj stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subObj); err != nil {
			break
		}
		ctx := r.Context()
		sub, err := s.db.Subscription().GetByStripeSubscriptionID(ctx, subObj.ID)
		if err != nil {
			break
		}
		// Update status based on Stripe subscription status
		switch subObj.Status {
		case "active":
			sub.Status = core.SubscriptionStatusActive
		case "trialing":
			sub.Status = core.SubscriptionStatusTrial
		case "canceled":
			sub.Status = core.SubscriptionStatusCanceled
		case "past_due":
			sub.Status = core.SubscriptionStatusPastDue
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
	default:
		// ignore other events
	}
	return s.renderJSON(w, http.StatusOK, statusResponse{Code: http.StatusOK, Message: "ok"})
}
