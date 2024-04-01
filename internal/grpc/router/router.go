package router

import (
	"context"
	"crypto/tls"
	"fmt"

	"gitlab.com/gtsh77-workshop/grpc-captcha/config"
	"gitlab.com/gtsh77-workshop/grpc-captcha/internal/grpc/controller/captcha"
	pb "gitlab.com/gtsh77-workshop/grpc-captcha/pkg/proto/grpc-captcha"
	"gitlab.com/gtsh77-workshop/grpc-captcha/pkg/tools"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	log *zap.SugaredLogger
	cfg *config.Config
}

func New(
	logger *zap.SugaredLogger,
	config *config.Config,
) *GrpcServer {
	return &GrpcServer{
		log: logger,
		cfg: config,
	}
}

func (s *GrpcServer) RegisterRouter() (*grpc.Server, error) {
	var (
		router    *grpc.Server
		tlsConfig *tls.Config
		err       error
	)

	if s.cfg.GRPC.TLS.Enabled {
		if tlsConfig, err = tools.PrepareTLS(s.cfg.GRPC.TLS.KeyData, s.cfg.GRPC.TLS.CrtData, s.cfg.GRPC.TLS.CrtCAData, s.cfg.GRPC.TLS.DomainName); err != nil {
			return nil, fmt.Errorf("tools.PrepareTLS: %w", err)
		}

		router = grpc.NewServer(
			grpc.ChainUnaryInterceptor(grpc_prometheus.UnaryServerInterceptor, s.localInterceptor),
			grpc.Creds(credentials.NewTLS(tlsConfig)))
	} else {
		router = grpc.NewServer(
			grpc.ChainUnaryInterceptor(grpc_prometheus.UnaryServerInterceptor, s.localInterceptor),
		)
	}

	pb.RegisterCaptchaServiceServer(router, &captcha.CaptchaService{
		Log: s.log,
		Cfg: s.cfg,
	})

	if s.cfg.GRPC.EnableProm {
		grpc_prometheus.Register(router)
	}

	return router, nil
}

func (s *GrpcServer) localInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var (
		ictx   context.Context
		cancel context.CancelFunc
		md     metadata.MD
		apiKey []string
		ok     bool

		m   interface{}
		err error
	)

	if md, ok = metadata.FromIncomingContext(ctx); ok {
		if apiKey, ok = md["x-api-key"]; ok && len(apiKey) == 1 {
			if apiKey[0] == s.cfg.GRPC.XApiKey {
				ictx, cancel = context.WithTimeout(ctx, s.cfg.GRPC.Timeout)
				defer cancel()

				if m, err = handler(ictx, req); err != nil {
					return m, err
				}

				return m, nil
			}
		}
	}

	return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
}
