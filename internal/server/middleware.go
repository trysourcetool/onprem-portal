package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/core"
	"github.com/trysourcetool/onprem-portal/internal/errdefs"
	"github.com/trysourcetool/onprem-portal/internal/jwt"
)

func (s *Server) authenticateUser(w http.ResponseWriter, r *http.Request) (*core.User, error) {
	ctx := r.Context()

	xsrfTokenHeader := r.Header.Get("X-XSRF-TOKEN")
	if xsrfTokenHeader == "" {
		return nil, errdefs.ErrUnauthenticated(errors.New("failed to get XSRF token"))
	}

	xsrfTokenCookie, err := r.Cookie("xsrf_token_same_site")
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	token, err := r.Cookie("access_token")
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	c, err := s.validateUserToken(token.Value)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	if err := validateXSRFToken(xsrfTokenHeader, xsrfTokenCookie.Value, c.XSRFToken); err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	userID, err := uuid.FromString(c.Subject)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	u, err := s.db.User().GetByID(ctx, userID)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	return u, nil
}

func (s *Server) validateUserToken(token string) (*jwt.AuthClaims, error) {
	if token == "" {
		return nil, errdefs.ErrUnauthenticated(errors.New("failed to get token"))
	}

	claims, err := jwt.ParseAuthClaims(token)
	if err != nil {
		return nil, errdefs.ErrUnauthenticated(err)
	}

	return claims, nil
}

func validateXSRFToken(header, cookie, claimToken string) error {
	if header == "" || cookie == "" || claimToken == "" {
		return errors.New("failed to get XSRF token")
	}
	if header != cookie && header != claimToken {
		return errors.New("invalid XSRF token")
	}
	return nil
}

func (s *Server) authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		u, err := s.authenticateUser(w, r)
		if err != nil {
			s.serveError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, internal.ContextUserKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
