package router

import (
	"gitlab.com/gtsh77-workshop/grpc-captcha/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func enableMiddleware(router *chi.Mux, cfg *config.Config) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Timeout(cfg.HTTP.Timeout))

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.HTTP.DomainNames,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
}
