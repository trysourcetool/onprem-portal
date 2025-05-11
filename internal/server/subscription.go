package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/core"
	"github.com/trysourcetool/onprem-portal/internal/errdefs"
)

type subscriptionResponse struct {
	ID                   string `json:"id"`
	UserID               string `json:"userId"`
	PlanID               string `json:"planId"`
	Status               string `json:"status"`
	StripeCustomerID     string `json:"stripeCustomerId"`
	StripeSubscriptionID string `json:"stripeSubscriptionId"`
	TrialStart           string `json:"trialStart"`
	TrialEnd             string `json:"trialEnd"`
	CreatedAt            string `json:"createdAt"`
	UpdatedAt            string `json:"updatedAt"`
}

func subscriptionFromModel(s *core.Subscription) *subscriptionResponse {
	if s == nil {
		return nil
	}
	return &subscriptionResponse{
		ID:                   s.ID.String(),
		UserID:               s.UserID.String(),
		PlanID:               s.PlanID.String(),
		Status:               s.Status.String(),
		StripeCustomerID:     s.StripeCustomerID,
		StripeSubscriptionID: s.StripeSubscriptionID,
		TrialStart:           strconv.FormatInt(s.TrialStart.Unix(), 10),
		TrialEnd:             strconv.FormatInt(s.TrialEnd.Unix(), 10),
		CreatedAt:            strconv.FormatInt(s.CreatedAt.Unix(), 10),
		UpdatedAt:            strconv.FormatInt(s.UpdatedAt.Unix(), 10),
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
	resp := &getSubscriptionResponse{Subscription: subscriptionFromModel(sub)}
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
	sub, err := s.db.Subscription().GetByUserID(ctx, ctxUser.ID)
	if err != nil {
		return err
	}
	sub.PlanID = planID
	sub.Status = core.SubscriptionStatusActive
	if err := s.db.Subscription().Update(ctx, sub); err != nil {
		return err
	}
	return s.renderJSON(w, http.StatusOK, statusResponse{Code: http.StatusOK, Message: "Subscription upgraded"})
}

func (s *Server) handleCancelSubscription(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	ctxUser := internal.ContextUser(ctx)
	sub, err := s.db.Subscription().GetByUserID(ctx, ctxUser.ID)
	if err != nil {
		return err
	}
	sub.Status = core.SubscriptionStatusCanceled
	if err := s.db.Subscription().Update(ctx, sub); err != nil {
		return err
	}
	return s.renderJSON(w, http.StatusOK, statusResponse{Code: http.StatusOK, Message: "Subscription canceled"})
}
