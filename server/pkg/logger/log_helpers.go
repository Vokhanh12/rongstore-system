package logger

import (
	"context"
	"os"
	"server/pkg/util/ctxutil"

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
		zap.String("trace_id", ctxutil.GetIDFromContext(ctx)),
		zap.String("request_id", ctxutil.RequestIdFromContext(ctx)),
		zap.String("status", p.Status),
		zap.Int64("latency_ms", p.LatencyMS),
		zap.Int("http_status", p.HTTPCode),
		zap.String("ip", p.IP),
		zap.String("user_agent", p.UserAgent),
		zap.String("instance", instance()),
		zap.String("env", env()),
		zap.String("event.type", "access"),
		zap.String("event.action", "request_handled"),
	}
	if sid := ctxutil.SessionIDFromContext(ctx); sid != "" {
		fields = append(fields, zap.String("session_id", sid))
	}
	if uid := ctxutil.UserIDFromContext(ctx); uid != "" {
		fields = append(fields, zap.String("user_id", uid))
	}
	for k, v := range p.Extra {
		fields = append(fields, zap.Any(k, v))
	}
	L.With(fields...).Info("request handled")
}

type AuditParams struct {
	Service string
	Action  string
	Msg     string
	Extra   map[string]interface{}
}

func LogAudit(ctx context.Context, p AuditParams) {
	fields := []zap.Field{
		zap.String("service", p.Service),
		zap.String("trace_id", ctxutil.GetIDFromContext(ctx)),
		zap.String("instance", instance()),
		zap.String("env", env()),
		zap.String("event.type", "audit"),
		zap.String("event.action", p.Action),
	}
	if sid := ctxutil.SessionIDFromContext(ctx); sid != "" {
		fields = append(fields, zap.String("session_id", sid))
	}
	if uid := ctxutil.UserIDFromContext(ctx); uid != "" {
		fields = append(fields, zap.String("user_id", uid))
	}
	for k, v := range p.Extra {
		fields = append(fields, zap.Any(k, v))
	}
	L.With(fields...).Info(p.Msg)
}

func LogSecurity(ctx context.Context, action, reason, msg string, extra map[string]interface{}) {
	fields := []zap.Field{
		zap.String("service", instance()),
		zap.String("trace_id", ctxutil.GetIDFromContext(ctx)),
		zap.String("instance", instance()),
		zap.String("env", env()),
		zap.String("event.type", "security"),
		zap.String("event.action", action),
		zap.String("security.reason", reason),
	}
	if sid := ctxutil.SessionIDFromContext(ctx); sid != "" {
		fields = append(fields, zap.String("session_id", sid))
	}
	if uid := ctxutil.UserIDFromContext(ctx); uid != "" {
		fields = append(fields, zap.String("user_id", uid))
	}
	for k, v := range extra {
		fields = append(fields, zap.Any(k, v))
	}
	L.With(fields...).Warn(msg)
}

func LogError(ctx context.Context, action, errMsg string, stack string, extra map[string]interface{}) {
	fields := []zap.Field{
		zap.String("service", instance()),
		zap.String("trace_id", ctxutil.GetIDFromContext(ctx)),
		zap.String("instance", instance()),
		zap.String("env", env()),
		zap.String("event.type", "error"),
		zap.String("event.action", action),
		zap.String("error.message", errMsg),
		zap.String("error.stack", stack),
	}
	if sid := ctxutil.SessionIDFromContext(ctx); sid != "" {
		fields = append(fields, zap.String("session_id", sid))
	}
	if uid := ctxutil.UserIDFromContext(ctx); uid != "" {
		fields = append(fields, zap.String("user_id", uid))
	}
	for k, v := range extra {
		fields = append(fields, zap.Any(k, v))
	}
	L.With(fields...).Error("error occurred")
}

func instance() string {
	if L == nil {
		return "unknown"
	}
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
