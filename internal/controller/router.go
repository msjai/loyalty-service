package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/msjai/loyalty-service/internal/config"
)

// NewRouter -.
func NewRouter(handler *chi.Mux, cfg *config.Config) *chi.Mux {
	handler.Use(middleware.Logger)

	return handler
}
