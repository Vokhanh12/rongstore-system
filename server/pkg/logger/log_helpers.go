package logger

import (
	"context"
	"os"

	"server/pkg/errors"

	"go.uber.org/zap"
)

type AccessParams struct {
	Service   string
	Handler   string
	Method    string
	HTTPCode  int
	Status    string
	LatencyMS int64
	IP        string
	UserAgent string
	Extra     map[string]interface{}
}

func LogAccess(ctx context.Context, p AccessParams) {
	fields := []zap.Field{
		zap.String("service", p.Service),
		zap.String("handler", p.Handler),
		zap.String("method", p.Method),
		zap.String("trace_id", errors.TraceIDFromContext(ctx)),
		zap.String("status", p.Status),
		zap.Int64("latency_ms", p.LatencyMS),
		zap.Int("http_status", p.HTTPCode),
		zap.String("ip", p.IP),
		zap.String("user_agent", p.UserAgent),
		zap.String("instance", instance()),
		zap.String("env", env()),
	}
	if sid := errors.SessionIDFromContext(ctx); sid != "" {
		fields = append(fields, zap.String("session_id", sid))
	}
	if uid := errors.UserIDFromContext(ctx); uid != "" {
		fields = append(fields, zap.String("user_id", uid))
	}
	for k, v := range p.Extra {
		fields = append(fields, zap.Any(k, v))
	}
	L.With(fields...).Info("request_handled")
}

type AuditParams struct {
	Service string
	Event   string
	Msg     string
	Extra   map[string]interface{}
}

func LogAudit(ctx context.Context, p AuditParams) {
	fields := []zap.Field{
		zap.String("service", p.Service),
		zap.String("event", p.Event),
		zap.String("trace_id", errors.TraceIDFromContext(ctx)),
		zap.String("instance", instance()),
		zap.String("env", env()),
	}
	if sid := errors.SessionIDFromContext(ctx); sid != "" {
		fields = append(fields, zap.String("session_id", sid))
	}
	if uid := errors.UserIDFromContext(ctx); uid != "" {
		fields = append(fields, zap.String("user_id", uid))
	}
	for k, v := range p.Extra {
		fields = append(fields, zap.Any(k, v))
	}
	// Audit goes to audit core (we used NewTee above to send all info-level to both access+audit sinks).
	L.With(fields...).Info(p.Msg)
}

func LogSecurity(ctx context.Context, event, reason, msg string, extra map[string]interface{}) {
	fields := []zap.Field{
		zap.String("service", instance()), // service can be passed as param instead
		zap.String("event", event),
		zap.String("reason", reason),
		zap.String("trace_id", errors.TraceIDFromContext(ctx)),
		zap.String("instance", instance()),
		zap.String("env", env()),
	}
	if sid := errors.SessionIDFromContext(ctx); sid != "" {
		fields = append(fields, zap.String("session_id", sid))
	}
	if uid := errors.UserIDFromContext(ctx); uid != "" {
		fields = append(fields, zap.String("user_id", uid))
	}
	for k, v := range extra {
		fields = append(fields, zap.Any(k, v))
	}
	L.With(fields...).Warn(msg)
}

func LogError(ctx context.Context, event, errMsg string, stack string, extra map[string]interface{}) {
	fields := []zap.Field{
		zap.String("service", instance()),
		zap.String("event", event),
		zap.String("trace_id", errors.TraceIDFromContext(ctx)),
		zap.String("instance", instance()),
		zap.String("env", env()),
		zap.String("error", errMsg),
		zap.String("stack", stack),
	}
	if sid := errors.SessionIDFromContext(ctx); sid != "" {
		fields = append(fields, zap.String("session_id", sid))
	}
	if uid := errors.UserIDFromContext(ctx); uid != "" {
		fields = append(fields, zap.String("user_id", uid))
	}
	for k, v := range extra {
		fields = append(fields, zap.Any(k, v))
	}
	L.With(fields...).Error("error_event")
}

// small helpers
func instance() string {
	if L == nil {
		return "unknown"
	}
	// cached hostname used in observability package — reuse env var if set
	if h := os.Getenv("INSTANCE"); h != "" {
		return h
	}
	if hn, err := os.Hostname(); err == nil {
		return hn
	}
	return "unknown"
}
func env() string {
	if e := os.Getenv("ENV"); e != "" {
		return e
	}
	return "unknown"
}
