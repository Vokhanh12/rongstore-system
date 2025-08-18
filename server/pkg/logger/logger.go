package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.Logger

func Init(prod bool) error {
	if prod {
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "ts"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		l, err := cfg.Build()
		if err != nil {
			return err
		}
		zap.ReplaceGlobals(l)
		L = l
		return nil
	}
	// dev
	cfg := zap.NewDevelopmentConfig()
	l, err := cfg.Build()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(l)
	L = l
	return nil
}
