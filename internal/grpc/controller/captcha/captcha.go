package captcha

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/redis/go-redis/v9"
	"gitlab.com/gtsh77-workshop/grpc-captcha/config"
	"gitlab.com/gtsh77-workshop/grpc-captcha/internal/grpc/controller"
	pb "gitlab.com/gtsh77-workshop/grpc-captcha/pkg/proto/grpc-captcha"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	NamespaceCaptcha = "grpc.controller.captcha"
)

type Controller struct {
	*controller.Base
	Rds *redis.Client
	pb.CaptchaServiceServer
}

func New(
	log *zap.SugaredLogger,
	cfg *config.Config,
	rds *redis.Client,
) *Controller {
	return &Controller{
		Base: controller.New(log, cfg),
		Rds:  rds,
	}
}

func (c *Controller) Generate(ctx context.Context, empt *empty.Empty) (*pb.CaptchaResponse, error) {
	var (
		id           string
		seq          []byte
		data, seqStr string
		err          error
	)

	id = uuid.NewString()

	seq = c.CaptchaNewSeq(c.Cfg.Render.DigCnt)
	seqStr = c.CaptchaSeqToString(seq)

	if err = c.Rds.Set(
		ctx,
		id,
		seqStr,
		c.Cfg.Render.TTL,
	).Err(); err != nil {
		c.Log.Named(NamespaceCaptcha).Errorf("redis.Set: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if data, err = c.CaptchaRenderImageHex(seq); err != nil {
		c.Log.Named(NamespaceCaptcha).Errorf("CaptchaRenderImageHex: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CaptchaResponse{
		Id:   id,
		Data: data,
	}, nil
}

func (c *Controller) Verify(ctx context.Context, req *pb.VerifyCaptchaRequest) (*empty.Empty, error) {
	var (
		rcmd *redis.StringCmd

		id, code, rdata string

		err error
	)

	if id = req.GetId(); len(id) == 0 {
		c.Log.Named(NamespaceCaptcha).Warn("CaptchaServiceServer.Verify empty field (id) provided")
		return nil, status.Error(codes.InvalidArgument, "empty field (id) provided")
	}

	if code = req.GetCode(); len(code) == 0 {
		c.Log.Named(NamespaceCaptcha).Warn("CaptchaServiceServer.Verify empty field (code) provided")
		return nil, status.Error(codes.InvalidArgument, "empty field (code) provided")
	}

	if rcmd = c.Rds.Get(ctx, id); err != nil {
		c.Log.Named(NamespaceCaptcha).Errorf("redis.Get: %v", err)
		return nil, status.Error(codes.Internal, err.Error())

	}

	if rdata, err = rcmd.Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			c.Log.Named(NamespaceCaptcha).Debugf("empty key request: %s", id)
			return nil, status.Error(codes.NotFound, codes.NotFound.String())
		}

		c.Log.Named(NamespaceCaptcha).Errorf("redis.Result: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if rdata != code {
		return nil, status.Error(codes.FailedPrecondition, codes.FailedPrecondition.String())
	}

	if err = c.Rds.Del(ctx, id).Err(); err != nil {
		c.Log.Named(NamespaceCaptcha).Warnf("redis.Del: %v", err)
	}

	return &empty.Empty{}, nil
}
