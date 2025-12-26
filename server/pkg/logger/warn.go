package logger

import "context"

type WarnParams struct {
	LogEntry
	Extra map[string]interface{}
}

func LogWarn(ctx context.Context, msg string, extra WarnParams) {
	fields := buildFieldsWarn(ctx, extra)
	WarnLogger.Warn(msg, fields...)
}
