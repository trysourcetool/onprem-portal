package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"strconv"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/config"
	"github.com/trysourcetool/onprem-portal/internal/core"
	"github.com/trysourcetool/onprem-portal/internal/database"
	"github.com/trysourcetool/onprem-portal/internal/errdefs"
	"github.com/trysourcetool/onprem-portal/internal/jwt"
	"github.com/trysourcetool/onprem-portal/internal/mail"
)

type userResponse struct {
	ID        string           `json:"id"`
	Email     string           `json:"email"`
	FirstName string           `json:"firstName"`
	LastName  string           `json:"lastName"`
	CreatedAt string           `json:"createdAt"`
	UpdatedAt string           `json:"updatedAt"`
	License   *licenseResponse `json:"license,omitempty"`
}

func (s *Server) userFromModel(user *core.User, l *core.License) *userResponse {
	if user == nil {
		return nil
	}

	return &userResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: strconv.FormatInt(user.CreatedAt.Unix(), 10),
		UpdatedAt: strconv.FormatInt(user.UpdatedAt.Unix(), 10),
		License:   s.licenseFromModel(l),
	}
}

type getMeResponse struct {
	User *userResponse `json:"user"`
}

func (s *Server) handleGetMe(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	ctxUser := internal.ContextUser(ctx)
	l, err := s.db.License().GetByUserID(ctx, ctxUser.ID)
	if err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, getMeResponse{
		User: s.userFromModel(ctxUser, l),
	})
}

type updateMeRequest struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

type updateMeResponse struct {
	User *userResponse `json:"user"`
}

func (s *Server) handleUpdateMe(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req updateMeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	ctxUser := internal.ContextUser(ctx)

	if req.FirstName != nil {
		ctxUser.FirstName = internal.StringValue(req.FirstName)
	}
	if req.LastName != nil {
		ctxUser.LastName = internal.StringValue(req.LastName)
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().Update(ctx, ctxUser); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, updateMeResponse{
		User: s.userFromModel(ctxUser, nil),
	})
}

type sendUpdateMeEmailInstructionsRequest struct {
	Email             string `json:"email" validate:"required,email"`
	EmailConfirmation string `json:"emailConfirmation" validate:"required,email"`
}

func buildUpdateEmailURL(token string) (string, error) {
	return internal.BuildURL(config.Config.BaseURL, path.Join("users", "email", "update", "confirm"), map[string]string{
		"token": token,
	})
}

func (s *Server) handleSendUpdateMeEmailInstructions(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req sendUpdateMeEmailInstructionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Validate email and confirmation match
	if req.Email != req.EmailConfirmation {
		return errdefs.ErrInvalidArgument(errors.New("email and email confirmation do not match"))
	}

	// Check if email already exists
	exists, err := s.db.User().IsEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return errdefs.ErrUserEmailAlreadyExists(errors.New("email already exists"))
	}

	// Get current user and organization
	ctxUser := internal.ContextUser(ctx)

	// Create token for email update
	tok, err := jwt.SignUpdateUserEmailToken(ctxUser.ID.String(), req.Email)
	if err != nil {
		return err
	}

	// Build update URL
	url, err := buildUpdateEmailURL(tok)
	if err != nil {
		return err
	}

	if err := mail.SendUpdateEmailInstructions(ctx, req.Email, ctxUser.FirstName, url); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, statusResponse{
		Code:    http.StatusOK,
		Message: "Email update instructions sent successfully",
	})
}

type updateMeEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type updateMeEmailResponse struct {
	User *userResponse `json:"user"`
}

func (s *Server) handleUpdateMeEmail(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req updateMeEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	c, err := jwt.ParseUpdateUserEmailClaims(req.Token)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	userID, err := uuid.FromString(c.Subject)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}
	u, err := s.db.User().GetByID(ctx, userID)
	if err != nil {
		return err
	}

	ctxUser := internal.ContextUser(ctx)
	if u.ID != ctxUser.ID {
		return errdefs.ErrUnauthenticated(errors.New("unauthorized"))
	}

	ctxUser.Email = c.Email

	if ctxUser.GoogleID != "" {
		ctxUser.GoogleID = ""
	}

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().Update(ctx, ctxUser); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, updateMeEmailResponse{
		User: s.userFromModel(ctxUser, nil),
	})
}
