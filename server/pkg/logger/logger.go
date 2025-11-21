package logger

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"

	"server/pkg/errors"
	"server/pkg/util/ctxutil"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// L is the global logger
var L *zap.Logger

// file cores
var (
	accessCore, auditCore, businessCore, infraCore, securityCore zapcore.Core
)

// -------------------- Severity → Zap Level mapping --------------------

var severityToZapLevel = map[string]zapcore.Level{
	"S1": zap.ErrorLevel, // Critical outage / security risk
	"S2": zap.WarnLevel,  // Degraded experience / important failure
	"S3": zap.InfoLevel,  // Minor / client input issue
}

func zapLevelFromSeverity(sev string) zapcore.Level {
	if lvl, ok := severityToZapLevel[sev]; ok {
		return lvl
	}
	return zap.InfoLevel // default
}

// -------------------- Init Logger --------------------

func Init() error {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)

	var err error
	accessCore, err = newFileCore("logs/access.log", jsonEncoder, zap.InfoLevel)
	if err != nil {
		return err
	}
	auditCore, err = newFileCore("logs/audit.log", jsonEncoder, zap.InfoLevel)
	if err != nil {
		return err
	}
	businessCore, err = newFileCore("logs/business_error.log", jsonEncoder, zap.WarnLevel)
	if err != nil {
		return err
	}
	infraCore, err = newFileCore("logs/infra_error.log", jsonEncoder, zap.ErrorLevel)
	if err != nil {
		return err
	}
	securityCore, err = newFileCore("logs/security.log", jsonEncoder, zap.WarnLevel)
	if err != nil {
		return err
	}

	core := zapcore.NewTee(accessCore, auditCore, businessCore, infraCore, securityCore)
	L = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(L)
	return nil

}

// helper to create file core
func newFileCore(path string, enc zapcore.Encoder, lvl zapcore.LevelEnabler) (zapcore.Core, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o640)
	if err != nil {
		return nil, err
	}
	return zapcore.NewCore(enc, zapcore.AddSync(f), lvl), nil
}

// -------------------- Helpers --------------------

func buildFields(ctx context.Context, extra map[string]interface{}) []zap.Field {
	fields := []zap.Field{
		zap.String("trace_id", ctxutil.GetIDFromContext(ctx)),
		zap.String("request_id", ctxutil.RequestIdFromContext(ctx)),
	}
	if sid := ctxutil.SessionIDFromContext(ctx); sid != "" {
		fields = append(fields, zap.String("session_id", sid))
	}
	if uid := ctxutil.UserIDFromContext(ctx); uid != "" {
		fields = append(fields, zap.String("user_id", uid))
	}
	fields = append(fields,
		zap.String("instance", instance()),
		zap.String("env", env()),
	)
	for k, v := range extra {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

func instance() string {
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

func TokenHashPrefix(token string) string {
	if token == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])[:8]
}

func MaskEmail(email string) string {
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			if i > 1 {
				return email[:1] + "***" + email[i:]
			}
			return "***" + email[i:]
		}
	}
	return "***"
}

// -------------------- Logging APIs --------------------

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
		zap.String("status", p.Status),
		zap.Int64("latency_ms", p.LatencyMS),
		zap.Int("http_status", p.HTTPCode),
		zap.String("ip", p.IP),
		zap.String("user_agent", p.UserAgent),
		zap.String("event.type", "access"),
		zap.String("event.action", "request_handled"),
	}
	fields = append(fields, buildFields(ctx, p.Extra)...)
	L.With(fields...).Info("access log")
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
		zap.String("event.type", "audit"),
		zap.String("event.action", p.Action),
	}
	fields = append(fields, buildFields(ctx, p.Extra)...)
	L.With(fields...).Info(p.Msg)
}

func LogSecurity(ctx context.Context, action, reason, msg string, extra map[string]interface{}) {
	fields := []zap.Field{
		zap.String("service", instance()),
		zap.String("event.type", "security"),
		zap.String("event.action", action),
		zap.String("security.reason", reason),
	}
	fields = append(fields, buildFields(ctx, extra)...)
	L.With(fields...).Warn(msg)
}

func LogBusinessError(ctx context.Context, err errors.BusinessError, extra map[string]interface{}) {
	fields := []zap.Field{
		zap.String("service", instance()),
		zap.String("event.type", "business_error"),
		zap.String("event.action", "business_failure"),
		zap.String("error.code", err.Code),
		zap.String("error.message", err.Message),
		zap.String("error.severity", err.Severity),
		zap.Bool("error.retryable", err.Retryable),
	}
	fields = append(fields, buildFields(ctx, extra)...)
	L.With(fields...).Warn("business error")
}

func LogInfraError(ctx context.Context, err errors.BusinessError, stack string, extra map[string]interface{}) {
	fields := []zap.Field{
		zap.String("service", instance()),
		zap.String("event.type", "infra_error"),
		zap.String("event.action", "infra_failure"),
		zap.String("error.code", err.Code),
		zap.String("error.message", err.Message),
		zap.String("error.severity", err.Severity),
		zap.Bool("error.retryable", err.Retryable),
		zap.String("error.stack", stack),
	}
	fields = append(fields, buildFields(ctx, extra)...)
	L.With(fields...).Error("infrastructure error")
}

func LogError(ctx context.Context, action, errMsg, stack string, extra map[string]interface{}) {
	fields := []zap.Field{
		zap.String("service", instance()),
		zap.String("event.type", "error"),
		zap.String("event.action", action),
		zap.String("error.message", errMsg),
		zap.String("error.stack", stack),
	}
	fields = append(fields, buildFields(ctx, extra)...)
	L.With(fields...).Error("error occurred")
}

// -------------------- Log by Severity --------------------

func LogBySeverity(ctx context.Context, err errors.BusinessError, extra map[string]interface{}) {
	level := zapLevelFromSeverity(err.Severity)
	switch level {
	case zap.ErrorLevel:
		LogInfraError(ctx, err, "", extra)
	case zap.WarnLevel:
		LogBusinessError(ctx, err, extra)
	default:
		LogAccess(ctx, AccessParams{Extra: extra})
	}
}
