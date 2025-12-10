package logger

import (
	"context"
)

func LogSecurity(ctx context.Context, msg string, extra map[string]interface{}) {
	fields := buildFields(ctx, extra)
	SecurityLogger.Warn(msg, fields...)
}
