package server

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/onprem-portal/internal/core"
	"github.com/trysourcetool/onprem-portal/internal/errdefs"
	"github.com/trysourcetool/onprem-portal/internal/jwt"
)

type refreshTokenResponse struct {
	ExpiresAt string `json:"expiresAt"`
}

func (s *Server) handleRefreshToken(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	xsrfTokenHeader := r.Header.Get("X-XSRF-TOKEN")
	if xsrfTokenHeader == "" {
		return errdefs.ErrUnauthenticated(errors.New("failed to get XSRF token"))
	}

	xsrfTokenCookie, err := r.Cookie("xsrf_token_same_site")
	if err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	// Validate XSRF token consistency
	if xsrfTokenCookie.Value != xsrfTokenHeader {
		return errdefs.ErrUnauthenticated(errors.New("invalid xsrf token"))
	}

	// Get user by refresh token
	hashedRefreshToken := core.HashRefreshToken(refreshTokenCookie.Value)
	u, err := s.db.User().GetByRefreshTokenHash(ctx, hashedRefreshToken)
	if err != nil {
		return errdefs.ErrUnauthenticated(err)
	}

	// Generate token and set expiration
	now := time.Now()
	expiresAt := now.Add(core.TokenExpiration())
	xsrfToken := uuid.Must(uuid.NewV4()).String()
	token, err := jwt.SignAuthToken(u.ID.String(), xsrfToken, expiresAt)
	if err != nil {
		return errdefs.ErrInternal(err)
	}

	cookieConfig := newCookieConfig()
	cookieConfig.SetAuthCookie(w, token, refreshTokenCookie.Value, xsrfToken,
		int(core.TokenExpiration().Seconds()),
		int(core.RefreshTokenExpiration.Seconds()),
		int(core.XSRFTokenExpiration.Seconds()),
	)

	return s.renderJSON(w, http.StatusOK, &refreshTokenResponse{
		ExpiresAt: strconv.FormatInt(expiresAt.Unix(), 10),
	})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) error {
	cookieConfig := newCookieConfig()
	cookieConfig.DeleteAuthCookie(w, r)

	return s.renderJSON(w, http.StatusOK, &statusResponse{
		Code:    http.StatusOK,
		Message: "Successfully logged out",
	})
}
