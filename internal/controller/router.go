package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/msjai/loyalty-service/internal/config"
	"github.com/msjai/loyalty-service/internal/usecase"
)

// NewRouter -.
func NewRouter(router *chi.Mux, loyalty usecase.Loyalty, cfg *config.Config) *chi.Mux {
	router.Use(middleware.Logger)

	// Routers
	router = newLoyaltyRoutes(router, loyalty, cfg)

	return router
}
