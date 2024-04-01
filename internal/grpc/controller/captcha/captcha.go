package captcha

import (
	"gitlab.com/gtsh77-workshop/grpc-captcha/config"
	pb "gitlab.com/gtsh77-workshop/grpc-captcha/pkg/proto/grpc-captcha"

	"go.uber.org/zap"
)

type CaptchaService struct {
	Log *zap.SugaredLogger
	Cfg *config.Config

	pb.CaptchaServiceServer
}
