package logger

import "context"

type InfoParams struct {
	LogEntry
}

func LogInfo(ctx context.Context, msg string, extra AccessParams) {
	fields := buildFieldsAccess(ctx, extra)
	AccessLogger.Info(msg, fields...)
}
