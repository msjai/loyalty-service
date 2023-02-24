package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/msjai/loyalty-service/internal/config"
	"github.com/msjai/loyalty-service/internal/controller"
	"github.com/msjai/loyalty-service/internal/usecase"
	"github.com/msjai/loyalty-service/internal/usecase/repo"
	"github.com/msjai/loyalty-service/internal/usecase/repo/postgres"
	"github.com/msjai/loyalty-service/internal/usecase/webapi"
)

// Run -.
func Run(cfg *config.Config) {
	l := cfg.L
	l.Infow("starting server...")

	repoPG, err := postgres.New(cfg.DataBaseURI, cfg.L)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - repo.New: %w", err))
	}
	defer repoPG.Close()

	// Use case
	loyaltyUseCase := usecase.New(
		repo.New(repoPG),
		webapi.New(cfg),
		cfg,
	)

	// initialize chi Mux object
	handler := chi.NewRouter()
	controller.NewRouter(handler, loyaltyUseCase, cfg)
	server := &http.Server{
		Addr:              cfg.RunAddress,
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}

	// Graceful server shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		l.Infof("Listening on port %v", server.Addr)
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCancel()

		l.Infow("Shutting down server...")
		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				l.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			l.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		l.Fatalf("listen: %s\n", err)
	}
	l.Infow("Server exiting")

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
