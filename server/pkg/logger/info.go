package logger

import "context"

type InfoParams struct {
	BaseLogLevel
}

func LogInfo(ctx context.Context, msg string, extra AccessParams) {
	fields := buildFieldsAccess(ctx, extra)
	AccessLogger.Info(msg, fields...)
}
