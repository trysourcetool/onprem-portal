package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid/v5"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/subscription"
	"github.com/stripe/stripe-go/v82/subscriptionitem"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/config"
	"github.com/trysourcetool/onprem-portal/internal/core"
	"github.com/trysourcetool/onprem-portal/internal/errdefs"
)

type subscriptionResponse struct {
	ID         string        `json:"id"`
	UserID     string        `json:"userId"`
	PlanID     string        `json:"planId"`
	Status     string        `json:"status"`
	TrialStart string        `json:"trialStart"`
	TrialEnd   string        `json:"trialEnd"`
	CreatedAt  string        `json:"createdAt"`
	UpdatedAt  string        `json:"updatedAt"`
	Plan       *planResponse `json:"plan"`
}

func subscriptionFromModel(sub *core.Subscription, plan *core.Plan) *subscriptionResponse {
	if sub == nil {
		return nil
	}
	return &subscriptionResponse{
		ID:     sub.ID.String(),
		UserID: sub.UserID.String(),
		PlanID: func() string {
			if sub.PlanID != nil {
				return sub.PlanID.String()
			}
			return ""
		}(),
		Plan: func() *planResponse {
			if plan != nil {
				return planFromModel(plan)
			}
			return nil
		}(),
		Status:     sub.Status.String(),
		TrialStart: strconv.FormatInt(sub.TrialStart.Unix(), 10),
		TrialEnd:   strconv.FormatInt(sub.TrialEnd.Unix(), 10),
		CreatedAt:  strconv.FormatInt(sub.CreatedAt.Unix(), 10),
		UpdatedAt:  strconv.FormatInt(sub.UpdatedAt.Unix(), 10),
	}
}

type getSubscriptionResponse struct {
	Subscription *subscriptionResponse `json:"subscription"`
}

func (s *Server) handleGetSubscription(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctxUser := internal.ContextUser(ctx)
	sub, err := s.db.Subscription().GetByUserID(ctx, ctxUser.ID)
	if err != nil {
		return err
	}
	var plan *core.Plan
	if sub.PlanID != nil {
		plan, err = s.db.Plan().GetByID(ctx, *sub.PlanID)
		if err != nil {
			return err
		}
	}
	resp := &getSubscriptionResponse{Subscription: subscriptionFromModel(sub, plan)}
	return s.renderJSON(w, http.StatusOK, resp)
}

type upgradeSubscriptionRequest struct {
	PlanID string `json:"planId"`
}

func (s *Server) handleUpgradeSubscription(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctxUser := internal.ContextUser(ctx)
	var req upgradeSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}
	planID, err := uuid.FromString(req.PlanID)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}
	plan, err := s.db.Plan().GetByID(ctx, planID)
	if err != nil {
		return err
	}
	sub, err := s.db.Subscription().GetByUserID(ctx, ctxUser.ID)
	if err != nil {
		return err
	}

	if sub.StripeSubscriptionID != "" {
		stripe.Key = config.Config.Stripe.Key
		stripeSub, err := subscription.Get(sub.StripeSubscriptionID, nil)
		if err != nil {
			return errdefs.ErrInternal(err)
		}
		if stripeSub.Items == nil || len(stripeSub.Items.Data) == 0 {
			return errdefs.ErrInternal(fmt.Errorf("stripe subscription items not found for user"))
		}
		params := &stripe.SubscriptionItemParams{
			Price:    stripe.String(plan.StripePriceID),
			Quantity: stripe.Int64(sub.SeatCount),
		}
		_, err = subscriptionitem.Update(stripeSub.Items.Data[0].ID, params)
		if err != nil {
			return errdefs.ErrInternal(err)
		}
	}

	sub.PlanID = &planID
	sub.Status = core.SubscriptionStatusActive
	if err := s.db.Subscription().Update(ctx, sub); err != nil {
		return err
	}
	return s.renderJSON(w, http.StatusOK, statusResponse{Code: http.StatusOK, Message: "Subscription upgraded"})
}
