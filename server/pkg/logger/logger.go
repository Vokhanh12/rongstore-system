package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	AccessLogger     *zap.Logger
	AuditLogger      *zap.Logger
	BusinessLogger   *zap.Logger
	InfraLogger      *zap.Logger
	DebugInfraLogger *zap.Logger
	SecurityLogger   *zap.Logger
)

func Init() error {
	encoderCfg := zapcore.EncoderConfig{TimeKey: "ts", LevelKey: "level", NameKey: "logger", CallerKey: "caller", MessageKey: "msg", StacktraceKey: "stack", LineEnding: zapcore.DefaultLineEnding, EncodeLevel: zapcore.LowercaseLevelEncoder, EncodeTime: zapcore.ISO8601TimeEncoder, EncodeDuration: zapcore.SecondsDurationEncoder, EncodeCaller: zapcore.ShortCallerEncoder}
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	var err error

	AccessLogger, err = newFileLogger("logs/access.log", encoder, zap.InfoLevel)
	if err != nil {
		return err
	}

	AuditLogger, err = newFileLogger("logs/audit.log", encoder, zap.InfoLevel)
	if err != nil {
		return err
	}

	BusinessLogger, err = newFileLogger("logs/business.log", encoder, zap.WarnLevel)
	if err != nil {
		return err
	}

	InfraLogger, err = newFileLogger("logs/infra.log", encoder, zap.ErrorLevel)
	if err != nil {
		return err
	}

	DebugInfraLogger, err = newFileLogger("logs/infra_debug.log", encoder, zap.DebugLevel)
	if err != nil {
		return err
	}

	SecurityLogger, err = newFileLogger("logs/security.log", encoder, zap.WarnLevel)
	if err != nil {
		return err
	}

	return nil
}
