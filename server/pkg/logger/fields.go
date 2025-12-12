package logger

import (
	"context"

	"go.uber.org/zap"
)

func buildFieldsAccess(ctx context.Context, p AccessParams) []zap.Field {
	fields := []zap.Field{
		zap.String("path", p.path),
		zap.String("method", p.method),
		zap.Int("http_code", p.httpCode),
		zap.String("ip", p.ip),
		zap.String("user_agent", p.userAgent),
	}

	for k, v := range p.extra {
		fields = append(fields, zap.Any(k, v))
	}

	return fields
}

func buildFields(ctx context.Context, extra map[string]interface{}) []zap.Field {
	var fields []zap.Field

	for k, v := range extra {
		fields = append(fields, zap.Any(k, v))
	}

	return fields
}
