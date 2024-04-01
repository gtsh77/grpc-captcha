package logger

import (
	"fmt"
	"net"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	Level         zapcore.Level
	UDPAddress    string
	IsJSONEncoder bool
	Color         bool
}

func NewLogger(config *LoggerConfig) (*zap.SugaredLogger, error) {
	var (
		encoderConfig zapcore.EncoderConfig
		encoder       zapcore.Encoder
		stdCore, core zapcore.Core
		cores         []zapcore.Core
		writer        zapcore.WriteSyncer
		conn          net.Conn
		logger        *zap.Logger

		err error
	)

	encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.ConsoleSeparator = " "

	if config.IsJSONEncoder {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	stdCore = zapcore.NewCore(encoder, os.Stdout, config.Level)

	cores = make([]zapcore.Core, 0)
	cores = append(cores, stdCore)

	if config.UDPAddress != "" {
		conn, err = net.Dial("udp", config.UDPAddress)
		if err != nil {
			return nil, fmt.Errorf("net.Dial: %w", err)
		}

		writer = zapcore.AddSync(conn)
		cores = append(cores, zapcore.NewCore(encoder, writer, config.Level))
	}

	core = zapcore.NewTee(cores...)
	logger = zap.New(core, zap.AddCaller())

	return logger.Sugar(), nil
}
