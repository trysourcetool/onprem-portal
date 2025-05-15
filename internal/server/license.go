package server

import (
	"errors"
	"net/http"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/core"
	"github.com/trysourcetool/onprem-portal/internal/errdefs"
)

type licenseResponse struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Key    string `json:"key"`
}

func (s *Server) licenseFromModel(l *core.License) *licenseResponse {
	if l == nil {
		return nil
	}

	key, err := s.encryptor.Decrypt(l.KeyCiphertext, l.KeyNonce)
	if err != nil {
		return nil
	}

	return &licenseResponse{
		ID:     l.ID.String(),
		UserID: l.UserID.String(),
		Key:    string(key),
	}
}

type licenseValidityResponse struct {
	Valid        bool                  `json:"valid"`
	Status       string                `json:"status"`
	Subscription *subscriptionResponse `json:"subscription,omitempty"`
}

func (s *Server) handleValidateLicense(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	license := internal.ContextLicense(ctx)
	if license == nil {
		return errdefs.ErrUnauthenticated(errors.New("license not found"))
	}
	sub, err := s.db.Subscription().GetByUserID(ctx, license.UserID)
	if err != nil {
		return errdefs.ErrLicenseNotFound(err)
	}
	var plan *core.Plan
	if sub.PlanID != nil {
		plan, err = s.db.Plan().GetByID(ctx, *sub.PlanID)
		if err != nil {
			return err
		}
	}
	valid := sub.IsActive() || sub.IsTrial()
	resp := &licenseValidityResponse{
		Valid:        valid,
		Status:       sub.Status.String(),
		Subscription: subscriptionFromModel(sub, plan),
	}
	return s.renderJSON(w, http.StatusOK, resp)
}
