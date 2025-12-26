package logger

import "context"

type ErrorParams struct {
	LogEntry LogEntry
	Extra    map[string]interface{}
}

func LogError(ctx context.Context, msg string, extra ErrorParams) {
	fields := buildFieldsError(ctx, extra)
	ErrorLogger.Error(msg, fields...)
}
