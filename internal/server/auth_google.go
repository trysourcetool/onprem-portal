package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/onprem-portal/internal/core"
	"github.com/trysourcetool/onprem-portal/internal/database"
	"github.com/trysourcetool/onprem-portal/internal/errdefs"
	"github.com/trysourcetool/onprem-portal/internal/google"
	"github.com/trysourcetool/onprem-portal/internal/jwt"
)

type requestGoogleAuthLinkResponse struct {
	AuthURL string `json:"authUrl"`
}

func (s *Server) handleRequestGoogleAuthLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	stateToken, err := jwt.SignGoogleAuthLinkToken()
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	googleOAuthClient := google.NewOAuthClient()
	url, err := googleOAuthClient.GetGoogleAuthCodeURL(ctx, stateToken)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	return s.renderJSON(w, http.StatusOK, &requestGoogleAuthLinkResponse{
		AuthURL: url,
	})
}

type authenticateWithGoogleRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

type authenticateWithGoogleResponse struct {
	ExpiresAt    string `json:"expiresAt"`
	Registration string `json:"registrationToken"`
	IsNewUser    bool   `json:"isNewUser"`
}

func (s *Server) handleAuthenticateWithGoogle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req authenticateWithGoogleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Parse and validate state token
	_, err := jwt.ParseGoogleAuthLinkClaims(req.State)
	if err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	// Get Google token and user info
	googleOAuthClient := google.NewOAuthClient()
	tok, err := googleOAuthClient.GetGoogleToken(ctx, req.Code)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	userInfo, err := googleOAuthClient.GetGoogleUserInfo(ctx, tok)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	// Check if user exists
	exists, err := s.db.User().IsEmailExists(ctx, userInfo.Email)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	if !exists {
		registrationToken, err := jwt.SignGoogleRegistrationToken(
			userInfo.ID,
			userInfo.Email,
			userInfo.GivenName,
			userInfo.FamilyName,
		)
		if err != nil {
			return errdefs.ErrInternal(fmt.Errorf("failed to create registration token: %w", err))
		}

		return s.renderJSON(w, http.StatusOK, &authenticateWithGoogleResponse{
			Registration: registrationToken,
			IsNewUser:    true,
		})
	}

	// For existing users
	u, err := s.db.User().GetByEmail(ctx, userInfo.Email)
	if err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	needsGoogleIDUpdate := u.GoogleID == ""

	xsrfToken := uuid.Must(uuid.NewV4()).String()
	now := time.Now()
	expiresAt := now.Add(core.TokenExpiration())

	token, err := jwt.SignAuthToken(u.ID.String(), xsrfToken, expiresAt)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	plainRefreshToken, hashedRefreshToken, err := core.GenerateRefreshToken()
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	u.RefreshTokenHash = hashedRefreshToken
	if needsGoogleIDUpdate {
		u.GoogleID = userInfo.ID
	}

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

	return s.renderJSON(w, http.StatusOK, &authenticateWithGoogleResponse{
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
		IsNewUser: false,
	})
}

type registerWithGoogleRequest struct {
	Token string `json:"token" validate:"required"`
}

type registerWithGoogleResponse struct {
	ExpiresAt string `json:"expiresAt"`
}

func (s *Server) handleRegisterWithGoogle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req registerWithGoogleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	if err := validateRequest(req); err != nil {
		return err
	}

	// Parse and validate registration token
	claims, err := jwt.ParseGoogleRegistrationClaims(req.Token)
	if err != nil {
		return errdefs.ErrInvalidArgument(fmt.Errorf("invalid registration token: %w", err))
	}

	// Check if user already exists
	exists, err := s.db.User().IsEmailExists(ctx, claims.Subject)
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to check user existence: %w", err))
	}
	if exists {
		return errdefs.ErrUserEmailAlreadyExists(fmt.Errorf("user with email %s already exists", claims.Subject))
	}

	plainRefreshToken, hashedRefreshToken, err := core.GenerateRefreshToken()
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to generate refresh token: %w", err))
	}

	tokenExpiration := core.TokenExpiration()
	u := &core.User{
		ID:               uuid.Must(uuid.NewV4()),
		Email:            claims.Subject,
		FirstName:        claims.FirstName,
		LastName:         claims.LastName,
		RefreshTokenHash: hashedRefreshToken,
		GoogleID:         claims.GoogleID,
	}

	plainLicenseKey, hashedLicenseKey, err := core.GenerateLicenseKey()
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to generate license key: %w", err))
	}

	ciphertext, nonce, err := s.encryptor.Encrypt([]byte(plainLicenseKey))
	if err != nil {
		return errdefs.ErrInternal(fmt.Errorf("failed to encrypt license key: %w", err))
	}

	l := &core.License{
		ID:            uuid.Must(uuid.NewV4()),
		UserID:        u.ID,
		KeyHash:       hashedLicenseKey,
		KeyCiphertext: ciphertext,
		KeyNonce:      nonce,
	}

	now := time.Now()
	expiresAt := now.Add(tokenExpiration)
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	var token string
	if err := s.db.WithTx(ctx, func(tx database.Tx) error {
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
			PlanID:     uuid.Nil, // Set to a default/free plan if needed
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

	return s.renderJSON(w, http.StatusOK, &registerWithGoogleResponse{
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
	})
}
