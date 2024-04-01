package router

import (
	"gitlab.com/gtsh77-workshop/grpc-captcha/config"
	"gitlab.com/gtsh77-workshop/grpc-captcha/internal/http/controller/health"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func New(
	log *zap.SugaredLogger,
	cfg *config.Config,
	rds *redis.Client,
	isRdy *bool,
) *chi.Mux {
	var (
		r  *chi.Mux
		ch *health.Controller
	)

	r = chi.NewRouter()
	enableMiddleware(r, cfg)

	if cfg.HTTP.EnablePprof {
		r.Mount("/debug", middleware.Profiler())
	}

	if cfg.HTTP.EnableProm {
		r.Handle(cfg.HTTP.MetricPath, promhttp.Handler())
	}

	ch = health.New(log, cfg, isRdy)

	r.Get(cfg.HTTP.HealthPath, ch.GetHealth)
	r.Get(cfg.HTTP.ReadyPath, ch.GetRdy)

	return r
}
