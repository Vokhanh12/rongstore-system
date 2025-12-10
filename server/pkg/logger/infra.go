package logger

import (
	"context"

	"go.uber.org/zap"
)

func LogInfraError(ctx context.Context, err error, extra map[string]interface{}) {
	fields := append([]zap.Field{zap.Error(err)}, buildFields(ctx, extra)...)
	InfraLogger.Error("infra error", fields...)
}

func LogInfraDebug(ctx context.Context, msg string, extra map[string]interface{}) {
	fields := buildFields(ctx, extra)
	DebugInfraLogger.Debug(msg, fields...)
}
