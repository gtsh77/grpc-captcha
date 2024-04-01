package controller

import (
	"gitlab.com/gtsh77-workshop/grpc-captcha/config"

	"go.uber.org/zap"
)

// all Controller specific deps (eg logger, config, validator)
type Base struct {
	Log *zap.SugaredLogger
	Cfg *config.Config
}

func New(
	log *zap.SugaredLogger,
	cfg *config.Config,
) *Base {
	return &Base{
		Log: log,
		Cfg: cfg,
	}
}
