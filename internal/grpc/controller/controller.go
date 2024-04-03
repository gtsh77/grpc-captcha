package controller

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/gtsh77-workshop/grpc-captcha/config"
	captcha "gitlab.com/gtsh77-workshop/grpc-captcha/pkg/dc-captcha-lite"

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

func (c *Base) CaptchaNewSeq(dcnt int) []byte {
	return captcha.RandomDigits(dcnt)
}

func (c *Base) CaptchaSeqToString(seq []byte) string {
	var (
		sb strings.Builder
		v  byte
	)

	for _, v = range seq {
		sb.WriteString(strconv.Itoa(int(v)))
	}

	return sb.String()
}

func (c *Base) CaptchaRenderImageHex(seq []byte) (string, error) {
	var (
		bb  *bytes.Buffer
		err error
	)

	bb = new(bytes.Buffer)
	if _, err = captcha.NewImage("", seq, c.Cfg.Render.Width, c.Cfg.Render.Height).WriteTo(bb); err != nil {
		return "", fmt.Errorf("captcha.Image.WriteTo: [%w]", err)
	}

	return base64.StdEncoding.EncodeToString(bb.Bytes()), nil
}
