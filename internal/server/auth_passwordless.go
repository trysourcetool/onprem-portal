package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/config"
	"github.com/trysourcetool/onprem-portal/internal/core"
	"github.com/trysourcetool/onprem-portal/internal/database"
	"github.com/trysourcetool/onprem-portal/internal/errdefs"
	"github.com/trysourcetool/onprem-portal/internal/jwt"
	"github.com/trysourcetool/onprem-portal/internal/mail"
)

func buildMagicLinkURL(token string) (string, error) {
	return internal.BuildURL(config.Config.BaseURL, path.Join("auth", "magic", "authenticate"), map[string]string{
		"token": token,
	})
}

type requestMagicLinkRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type requestMagicLinkResponse struct {
	Email string `json:"email"`
	IsNew bool   `json:"isNew"`
}

func (s *Server) handleRequestMagicLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req requestMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Check if email exists
	exists, err := s.db.User().IsEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}

	var firstName string
	if exists {
		// Get user by email for existing users
		u, err := s.db.User().GetByEmail(ctx, req.Email)
		if err != nil {
			return err
		}
		firstName = u.FirstName
	} else {
		// For new users, generate a temporary ID that will be verified/used later
		firstName = "there" // Default greeting
	}

	// Create token for magic link authentication
	tok, err := jwt.SignMagicLinkToken(req.Email)
	if err != nil {
		return err
	}

	// Build magic link URL
	url, err := buildMagicLinkURL(tok)
	if err != nil {
		return err
	}

	// Send magic link email
	if err := mail.SendMagicLinkEmail(ctx, req.Email, firstName, url); err != nil {
		return err
	}

	return s.renderJSON(w, http.StatusOK, requestMagicLinkResponse{
		Email: req.Email,
		IsNew: !exists,
	})
}

type authenticateWithMagicLinkRequest struct {
	Token     string `json:"token" validate:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type authenticateWithMagicLinkResponse struct {
	RegistrationToken string `json:"registrationToken"`
	ExpiresAt         string `json:"expiresAt"`
	IsNewUser         bool   `json:"isNewUser"`
}

func (s *Server) handleAuthenticateWithMagicLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req authenticateWithMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	c, err := jwt.ParseMagicLinkClaims(req.Token)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	// Check if user exists
	exists, err := s.db.User().IsEmailExists(ctx, c.Subject)
	if err != nil {
		return err
	}

	if !exists {
		// Generate registration token for new user
		registrationToken, err := jwt.SignMagicLinkRegistrationToken(c.Subject)
		if err != nil {
			return fmt.Errorf("failed to generate registration token: %w", err)
		}

		return s.renderJSON(w, http.StatusOK, authenticateWithMagicLinkResponse{
			RegistrationToken: registrationToken,
			IsNewUser:         true,
		})
	}

	// Get existing user
	u, err := s.db.User().GetByEmail(ctx, c.Subject)
	if err != nil {
		return err
	}

	now := time.Now()
	expiresAt := now.Add(core.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()

	token, err := jwt.SignAuthToken(u.ID.String(), xsrfToken, expiresAt)
	if err != nil {
		return err
	}

	plainRefreshToken, hashedRefreshToken, err := core.GenerateRefreshToken()
	if err != nil {
		return err
	}

	u.RefreshTokenHash = hashedRefreshToken

	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		if err := tx.User().Update(ctx, u); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	cookieConfig := newCookieConfig()
	cookieConfig.SetAuthCookie(w, token, plainRefreshToken, xsrfToken,
		int(core.TokenExpiration().Seconds()),
		int(core.RefreshTokenExpiration.Seconds()),
		int(core.XSRFTokenExpiration.Seconds()),
	)

	return s.renderJSON(w, http.StatusOK, authenticateWithMagicLinkResponse{
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
		IsNewUser: false,
	})
}

type registerWithMagicLinkRequest struct {
	Token     string `json:"token" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type registerWithMagicLinkResponse struct {
	ExpiresAt string `json:"expiresAt"`
}

func (s *Server) handleRegisterWithMagicLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req registerWithMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Parse and validate the registration token
	claims, err := jwt.ParseMagicLinkRegistrationClaims(req.Token)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	// Generate refresh token and XSRF token
	plainRefreshToken, hashedRefreshToken, err := core.GenerateRefreshToken()
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	// Create a new user
	u := &core.User{
		ID:               uuid.Must(uuid.NewV4()),
		Email:            claims.Subject,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		RefreshTokenHash: hashedRefreshToken,
	}

	plainLicenseKey, hashedLicenseKey, err := core.GenerateLicenseKey()
	if err != nil {
		return err
	}

	ciphertext, nonce, err := s.encryptor.Encrypt([]byte(plainLicenseKey))
	if err != nil {
		return err
	}

	l := &core.License{
		ID:            uuid.Must(uuid.NewV4()),
		UserID:        u.ID,
		KeyHash:       hashedLicenseKey,
		KeyCiphertext: ciphertext,
		KeyNonce:      nonce,
	}

	xsrfToken := uuid.Must(uuid.NewV4()).String()
	now := time.Now()
	expiresAt := now.Add(core.TokenExpiration())

	var token string
	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
		// Create the user in a transaction
		if err := tx.User().Create(ctx, u); err != nil {
			return err
		}

		if err := tx.License().Create(ctx, l); err != nil {
			return err
		}

		// Create subscription (trial)
		trialStart := now
		trialEnd := now.Add(time.Duration(core.TrialPeriodDays) * 24 * time.Hour)
		sub := &core.Subscription{
			ID:         uuid.Must(uuid.NewV4()),
			UserID:     u.ID,
			PlanID:     uuid.Nil,
			Status:     core.SubscriptionStatusTrial,
			TrialStart: trialStart,
			TrialEnd:   trialEnd,
		}
		if err := tx.Subscription().Create(ctx, sub); err != nil {
			return err
		}

		token, err = jwt.SignAuthToken(u.ID.String(), xsrfToken, expiresAt)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	cookieConfig := newCookieConfig()
	cookieConfig.SetAuthCookie(w, token, plainRefreshToken, xsrfToken,
		int(core.TokenExpiration().Seconds()),
		int(core.RefreshTokenExpiration.Seconds()),
		int(core.XSRFTokenExpiration.Seconds()),
	)

	return s.renderJSON(w, http.StatusOK, registerWithMagicLinkResponse{
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
	})
}
