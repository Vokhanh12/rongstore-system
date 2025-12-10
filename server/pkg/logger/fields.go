package logger

import (
	"context"

	"go.uber.org/zap"
)

func buildFields(ctx context.Context, extra map[string]interface{}) []zap.Field {
	var fields []zap.Field

	for k, v := range extra {
		fields = append(fields, zap.Any(k, v))
	}

	return fields
}
