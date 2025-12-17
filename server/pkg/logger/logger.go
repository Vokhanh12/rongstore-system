package logger

import (
	"context"

	"server/pkg/errors"

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

type BaseLogging struct {
	Service string
}

type PayloadLogging struct {
	TraceId  string
	UserId   string
	ClientId string
	RealmId  string
}

type BaseLogLevel struct {
	BaseLog      BaseLogging = BASE_LOGGING
	PayloadLog   PayloadLogging
	Code         string `json:"code"`
	Key          string `json:"key"`
	expected     bool
	HTTPStatus   string `json:"http_status"`
	GRPCCode     string `json:"grpc_code"`
	Message      string `json:"message"`
	Cause        string `json:"cause"`
	CauseDetail  string `json:"cause_detail"`
	ClientAction string `json:"client_action"`
	ServerAction string `json:"server_action"`
}

func Init(service string) error {
	BASE_LOGGING.Service = service

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

func LogBySeverity(ctx context.Context, msg string, err *errors.AppError, extra map[string]interface{}) {

	level := LevelBySeverity(err.Severity, err.Expected)
	switch level {
	case zapcore.ErrorLevel:
		LogError(ctx, msg, ErrorParams{
			BaseLogLevel: BASE_LOGGING,
		})
	case zapcore.WarnLevel:
		logger.Warn(msg, fields...)
	case zapcore.InfoLevel:
		logger.Info(msg, fields...)
	case zapcore.DebugLevel:
		logger.Debug(msg, fields...)
	default:
		logger.Info(msg, fields...)
	}
}
