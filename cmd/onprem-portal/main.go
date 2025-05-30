package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/trysourcetool/onprem-portal/internal/config"
	"github.com/trysourcetool/onprem-portal/internal/encrypt"
	"github.com/trysourcetool/onprem-portal/internal/logger"
	"github.com/trysourcetool/onprem-portal/internal/postgres"
	"github.com/trysourcetool/onprem-portal/internal/server"
)

func init() {
	config.Init()
	logger.Init()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pqClient, err := postgres.Open()
	if err != nil {
		logger.Logger.Fatal("failed to open postgres", zap.Error(err))
	}

	db := postgres.New(pqClient)
	encryptor, err := encrypt.NewEncryptor()
	if err != nil {
		logger.Logger.Fatal("failed to create encryptor", zap.Error(err))
	}

	// if config.Config.Env == config.EnvLocal {
	// 	if err := internal.LoadFixtures(ctx, db); err != nil {
	// 		logger.Logger.Fatal(err.Error())
	// 	}
	// }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Logger.Info(fmt.Sprintf("Defaulting to port %s\n", port))
	}

	handler := chi.NewRouter()
	s := server.New(db, encryptor)
	s.Install(handler)

	srv := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      600 * time.Second,
		Handler:           handler,
		Addr:              fmt.Sprintf(":%s", port),
	}

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		logger.Logger.Info(fmt.Sprintf("Listening on port %s\n", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("HTTP server error: %v", err)
		}
		return nil
	})
	eg.Go(func() error {
		<-egCtx.Done()
		logger.Logger.Info("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		// Attempt to gracefully shut down the server first.
		var shutdownErr error
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Logger.Error("Server shutdown error", zap.Error(err))
			shutdownErr = fmt.Errorf("server shutdown: %v", err)
		}

		if err := pqClient.Close(); err != nil {
			logger.Logger.Sugar().Errorf("DB connection close failed: %v", err)
		} else {
			logger.Logger.Sugar().Info("DB connection gracefully stopped")
		}

		logger.Logger.Info("Server shutdown complete")
		// Return the server shutdown error if it happened.
		return shutdownErr
	})

	if err := eg.Wait(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Error(fmt.Sprintf("Error during shutdown: %v", err))
		os.Exit(1)
	}
}
