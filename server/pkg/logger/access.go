package logger

import (
	"context"
)

func LogAccess(ctx context.Context, msg string, extra map[string]interface{}) {
	fields := buildFields(ctx, extra)
	AccessLogger.Info(msg, fields...)
}
