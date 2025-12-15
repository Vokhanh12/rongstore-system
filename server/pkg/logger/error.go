package logger

import "context"

type ErrorParams struct {
	BaseLogLevel
}

func LogErr(ctx context.Context, msg string, extra AccessParams) {
	fields := buildFieldsAccess(ctx, extra)
	AccessLogger.Info(msg, fields...)
}
