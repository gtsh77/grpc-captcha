package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"gitlab.com/gtsh77-workshop/grpc-captcha/config"
	grouter "gitlab.com/gtsh77-workshop/grpc-captcha/internal/grpc/router"
	"gitlab.com/gtsh77-workshop/grpc-captcha/internal/http/router"
	"gitlab.com/gtsh77-workshop/grpc-captcha/pkg/logger"
	"gitlab.com/gtsh77-workshop/grpc-captcha/pkg/tools"
	"google.golang.org/grpc"

	"github.com/ardanlabs/conf/v3"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Service struct {
	log  *zap.SugaredLogger
	Cfg  *config.Config
	Rds  *redis.Client
	Pgs  *sqlx.DB
	http *http.Server
	rpc  *grpc.Server
	Rtr  *chi.Mux

	isRdy  bool
	stopCh chan bool
}

func New(name, version, compliedAt string) *Service {
	return &Service{
		Cfg:    config.New(name, version, compliedAt),
		stopCh: make(chan bool),
	}
}

func (s *Service) Start() (*Service, error) {
	var (
		ctx context.Context
		cf  context.CancelFunc
		err error
	)

	ctx, cf = context.WithCancel(context.Background())
	defer s.recover()
	defer cf()

	if err = s.readEnvAndArgs(); err != nil {
		return nil, fmt.Errorf("readEnvAndArgs: %w", err)
	}

	if err = s.startLogger(); err != nil {
		return nil, fmt.Errorf("startLogger: %w", err)
	}

	if s.Cfg.Runtime.ID, err = tools.RandHex(3); err != nil {
		s.log.Warnf("tools.RandHex: %v", err)
	}

	if err = s.connectRedis(ctx); err != nil {
		s.log.Fatalf("connectRedis: %v", err)
	}
	s.log.Infof("redis instance connected: %s:%s", s.Cfg.Redis.Host, s.Cfg.Redis.Port)

	s.Rtr = router.New(s.log, s.Cfg, s.Rds, &s.isRdy)

	sigCh := s.subSig()
	httpCh := s.startHTTP(ctx, s.Rtr)
	grpcCh := s.StartGrpc()

	go s.gsWatcher(ctx, sigCh, httpCh, grpcCh)
	s.log.Infof("HTTP server enabled: %s:%s", s.Cfg.HTTP.Host, s.Cfg.HTTP.Port)
	s.log.Infof("gRPC server enabled: %s:%s", s.Cfg.GRPC.Host, s.Cfg.GRPC.Port)

	s.isRdy = true

	if flag.Lookup("test.v") == nil {
		<-s.stopCh
	}

	return s, nil
}

func (s *Service) recover() {
	var (
		rec  interface{}
		buf  [2048]byte
		blen int
	)

	if rec = recover(); rec != nil {
		blen = runtime.Stack(buf[:], true)
		os.Stdout.Write(buf[:blen])
	}
}

func (s *Service) startLogger() error {
	var err error

	if s.log, err = logger.NewLogger(&logger.LoggerConfig{
		Level:         zapcore.Level(s.Cfg.Log.Level),
		IsJSONEncoder: s.Cfg.Log.AsJSON,
	}); err != nil {
		return err
	}

	return err
}

func (s *Service) readEnvAndArgs() error {
	var (
		h   string
		err error
	)

	if h, err = conf.Parse(config.ServicePrefix, s.Cfg); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(h) //nolint:forbidigo //non-json ouptput required
			return nil
		}

		return err
	}

	return nil
}

func (s *Service) subSig() <-chan os.Signal {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGSEGV,
		syscall.SIGBUS,
		syscall.SIGILL)

	return sigCh
}

func (s *Service) connectRedis(ctx context.Context) error {
	var (
		opts *redis.Options
		cmd  *redis.StatusCmd

		err error
	)

	opts = &redis.Options{
		Addr:            net.JoinHostPort(s.Cfg.Redis.Host, s.Cfg.Redis.Port),
		Username:        s.Cfg.Redis.User,
		Password:        s.Cfg.Redis.Pass,
		DB:              s.Cfg.Redis.DB,
		DialTimeout:     s.Cfg.Redis.MaxExecTime,
		ReadTimeout:     s.Cfg.Redis.MaxExecTime,
		WriteTimeout:    s.Cfg.Redis.MaxExecTime,
		MaxRetries:      s.Cfg.Redis.MaxRetries,
		MaxActiveConns:  s.Cfg.Redis.MaxOpenConn,
		MaxIdleConns:    s.Cfg.Redis.MaxIdleConn,
		ConnMaxLifetime: s.Cfg.Redis.MaxConnTTL,
	}

	if s.Cfg.Redis.TLS.Enabled {
		if opts.TLSConfig, err = tools.PrepareTLS(s.Cfg.Redis.TLS.KeyData, s.Cfg.Redis.TLS.CrtData, s.Cfg.Redis.TLS.CrtCAData, s.Cfg.Redis.TLS.DomainName); err != nil {
			return fmt.Errorf("tools.PrepareTLS: %w", err)
		}
	}

	s.Rds = redis.NewClient(opts)
	if cmd = s.Rds.Ping(ctx); cmd.Err() != nil {
		return err
	}

	return nil
}

func (s *Service) startHTTP(ctx context.Context, handler *chi.Mux) <-chan error {
	var (
		httpCh chan error
		err    error
	)

	s.http = &http.Server{
		ReadTimeout:  s.Cfg.HTTP.Timeout,
		WriteTimeout: s.Cfg.HTTP.Timeout,
		Addr:         net.JoinHostPort(s.Cfg.HTTP.Host, s.Cfg.HTTP.Port),
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
		Handler: handler,
	}

	go func() {
		httpCh = make(chan error, 1)
		defer close(httpCh)

		if s.Cfg.HTTP.TLS.Enabled {
			if s.http.TLSConfig, err = tools.PrepareTLS(s.Cfg.HTTP.TLS.KeyData, s.Cfg.HTTP.TLS.CrtData, s.Cfg.HTTP.TLS.CrtCAData, s.Cfg.HTTP.TLS.DomainName); err != nil {
				httpCh <- fmt.Errorf("tools.PrepareTLS: %w", err)
			}

			httpCh <- s.http.ListenAndServeTLS("", "")
		} else {
			httpCh <- s.http.ListenAndServe()
		}
	}()

	return httpCh
}

func (s *Service) StartGrpc() <-chan error {
	var (
		err      error
		listener net.Listener
		grpcCh   chan error
	)

	grpcCh = make(chan error, 1)

	if s.rpc, err = grouter.New(s.log, s.Cfg, s.Rds).RegisterRouter(); err != nil {
		grpcCh <- err
		return grpcCh
	}

	if listener, err = net.Listen("tcp", fmt.Sprint(s.Cfg.GRPC.Host, ":", s.Cfg.GRPC.Port)); err != nil {
		grpcCh <- err
		return grpcCh
	}

	go func() {
		defer close(grpcCh)
		grpcCh <- s.rpc.Serve(listener)
	}()

	return grpcCh
}

func (s *Service) gsWatcher(ctx context.Context, sigCh <-chan os.Signal, httpCh, grpcCh <-chan error) {
	select {
	case sig := <-sigCh:
		s.log.Warnf("os signal: %s", sig.String())
	case err := <-httpCh:
		s.log.Errorf("http.ListenAndServe: %v", err)
	case err := <-grpcCh:
		s.log.Errorf("grpc router: %v", err)
	}

	s.shutdown(ctx)
}

func (s *Service) shutdown(ctx context.Context) error {
	var err error

	s.isRdy = false

	s.rpc.GracefulStop()
	s.log.Warn("grpc server shutdown")

	if s.http != nil {
		if err = s.http.Shutdown(ctx); err != nil {
			s.log.Errorf("http.Shutdown: %v", err)
		} else {
			s.log.Warn("http server shutdown")
		}
	}

	if s.Rds != nil {
		if err := s.Rds.Close(); err != nil {
			s.log.Errorf("redis.Close: %v", err)
		} else {
			s.log.Warn("redis connection closed")
		}
	}

	s.stopCh <- true

	return nil
}
