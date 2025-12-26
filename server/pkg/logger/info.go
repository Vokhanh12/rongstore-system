package logger

import "context"

type InfoParams struct {
	LogEntry
	Extra map[string]interface{}
}

func LogInfo(ctx context.Context, msg string, extra InfoParams) {
	fields := buildFieldsInfo(ctx, extra)
	InfoLogger.Info(msg, fields...)
}
