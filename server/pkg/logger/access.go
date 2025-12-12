package logger

import (
	"context"
)

type AccessParams struct {
	Base      BaseLogger
	Path      string
	Method    string
	HTTPCode  int
	IP        string
	UserAgent string
	LatencyMS int64
	Extra     map[string]interface{}
}

func LogAccess(ctx context.Context, msg string, extra AccessParams) {
	fields := buildFieldsAccess(ctx, extra)
	AccessLogger.Info(msg, fields...)
}
