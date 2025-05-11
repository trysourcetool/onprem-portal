package server

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/config"
	"github.com/trysourcetool/onprem-portal/internal/database"
	"github.com/trysourcetool/onprem-portal/internal/encrypt"
	"github.com/trysourcetool/onprem-portal/internal/errdefs"
	"github.com/trysourcetool/onprem-portal/internal/logger"
)

type Server struct {
	db        database.DB
	encryptor *encrypt.Encryptor
}

func New(db database.DB, encryptor *encrypt.Encryptor) *Server {
	return &Server{db, encryptor}
}

func (s *Server) installDefaultMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(time.Duration(600) * time.Second))
}

func (s *Server) installCORSMiddleware(router *chi.Mux) {
	router.Use(cors.New(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			return true
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{},
		AllowCredentials: true,
		MaxAge:           0,
		Debug:            !(config.Config.Env == config.EnvProd),
	}).Handler)
}

func (s *Server) errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			s.serveError(w, r, err)
		}
	}
}

func (s *Server) installRESTHandlers(router *chi.Mux) {
	router.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		})

		r.Route("/v1", func(r chi.Router) {
			r.Route("/auth", func(r chi.Router) {
				r.Post("/magic/request", s.errorHandler(s.handleRequestMagicLink))
				r.Post("/magic/authenticate", s.errorHandler(s.handleAuthenticateWithMagicLink))
				r.Post("/magic/register", s.errorHandler(s.handleRegisterWithMagicLink))

				r.Post("/google/request", s.errorHandler(s.handleRequestGoogleAuthLink))
				r.Post("/google/authenticate", s.errorHandler(s.handleAuthenticateWithGoogle))
				r.Post("/google/register", s.errorHandler(s.handleRegisterWithGoogle))

				r.Post("/refreshToken", s.errorHandler(s.handleRefreshToken))
				r.Post("/logout", s.errorHandler(s.handleLogout))
			})

			r.Route("/users", func(r chi.Router) {
				r.Use(s.authUser)

				r.Route("/me", func(r chi.Router) {
					r.Get("/", s.errorHandler(s.handleGetMe))
					r.Put("/", s.errorHandler(s.handleUpdateMe))
					r.Post("/email/instructions", s.errorHandler(s.handleSendUpdateMeEmailInstructions))
					r.Put("/email", s.errorHandler(s.handleUpdateMeEmail))
				})
			})

			r.Route("/subscriptions", func(r chi.Router) {
				r.Use(s.authUser)
				r.Get("/", s.errorHandler(s.handleGetSubscription))
				r.Post("/upgrade", s.errorHandler(s.handleUpgradeSubscription))
				r.Post("/cancel", s.errorHandler(s.handleCancelSubscription))
			})

			r.Get("/plans", s.errorHandler(s.handleListPlans))
		})
	})
}

func (s *Server) installStaticHandler(router *chi.Mux) {
	staticDir := os.Getenv("STATIC_FILES_DIR")
	serveStaticFiles(router, staticDir)
}

func (s *Server) Install(router *chi.Mux) {
	s.installDefaultMiddlewares(router)
	s.installCORSMiddleware(router)
	s.installRESTHandlers(router)
	s.installStaticHandler(router)
}

func (s *Server) serveError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	ctxUser := internal.ContextUser(ctx)
	var email string
	if ctxUser != nil {
		email = ctxUser.Email
	}

	v, ok := err.(*errdefs.Error)
	if !ok {
		logger.Logger.Error(
			err.Error(),
			zap.Stack("stack_trace"),
			zap.String("email", email),
			zap.String("cause", "application"),
		)

		s.renderJSON(
			w,
			http.StatusInternalServerError,
			errdefs.ErrInternal(err),
		)
		return
	}

	fields := []zap.Field{
		zap.String("email", email),
		zap.String("error_stacktrace", strings.Join(v.StackTrace(), "\n")),
	}

	switch {
	case v.Status >= 500:
		fields = append(fields, zap.String("cause", "application"))
		logger.Logger.Error(err.Error(), fields...)
	case v.Status >= 402, v.Status == 400:
		fields = append(fields, zap.String("cause", "user"))
		logger.Logger.Error(err.Error(), fields...)
	default:
		fields = append(fields, zap.String("cause", "internal_info"))
		logger.Logger.Warn(err.Error(), fields...)
	}

	s.renderJSON(w, v.Status, v)
}

func (s *Server) renderJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

type statusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
