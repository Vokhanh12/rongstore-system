package logger

import (
	"context"
)

func LogAudit(ctx context.Context, msg string, extra map[string]interface{}) {
	fields := buildFields(ctx, extra)
	AuditLogger.Info(msg, fields...)
}
