package health

import (
	"encoding/json"
	"net/http"

	"gitlab.com/gtsh77-workshop/grpc-captcha/config"
	"gitlab.com/gtsh77-workshop/grpc-captcha/internal/entity"
	"gitlab.com/gtsh77-workshop/grpc-captcha/internal/http/controller"

	"go.uber.org/zap"
)

type Controller struct {
	*controller.Base
	isRdy *bool
}

func New(
	log *zap.SugaredLogger,
	cfg *config.Config,
	isRdy *bool,
) *Controller {
	return &Controller{
		Base:  controller.New(log, cfg),
		isRdy: isRdy,
	}
}

func (c *Controller) GetHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&entity.HeathCheckResponse{
		ID:         c.Cfg.Runtime.ID,
		Name:       c.Cfg.Runtime.Name,
		Version:    c.Cfg.Runtime.Version,
		CompiledAt: c.Cfg.Runtime.CompiledAt,
	}); err != nil {
		c.Log.Errorf("json.Encoder: %v", err)
	}
}

func (c *Controller) GetRdy(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if c.isRdy == nil || !*c.isRdy {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(&entity.HealthOperableResponse{
		IsRdy: *c.isRdy,
	}); err != nil {
		c.Log.Errorf("json.Encoder: %v", err)
	}
}
