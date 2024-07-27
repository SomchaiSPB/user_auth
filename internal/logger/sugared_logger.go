package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewSugaredLogger() (*zap.SugaredLogger, error) {
	var cfg zap.Config

	cfg = zap.NewDevelopmentConfig()

	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	l, err := cfg.Build()

	if err != nil {
		return nil, err
	}

	return l.Sugar(), nil
}
