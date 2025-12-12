package ctxutil

import "context"

type ctxKey string

const (
	ctxKeyTraceID   ctxKey = "trace_id"
	ctxKeySessionID ctxKey = "session_id"
	ctxKeyUserID    ctxKey = "user_id"
	ctxKeyRequestID ctxKey = "request_id"
	ctxKeyClientId  ctxKey = "client_id"
	ctxKeyRealmId   ctxKey = "realm_id"
)

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ctxKeyTraceID, traceID)
}

func GetIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyTraceID).(string); ok {
		return v
	}
	return ""
}

func RequestIdFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyRequestID).(string); ok {
		return v
	}
	return ""
}

func TraceIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyTraceID).(string); ok {
		return v
	}
	return ""
}

func ClientIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyTraceID).(string); ok {
		return v
	}
	return ""
}

func RealmIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyTraceID).(string); ok {
		return v
	}
	return ""
}

func WithSessionID(ctx context.Context, sid string) context.Context {
	return context.WithValue(ctx, ctxKeySessionID, sid)
}

func WithRequestID(ctx context.Context, sid string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, sid)
}

func SessionIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeySessionID).(string); ok {
		return v
	}
	return ""
}

func WithUserID(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, ctxKeyUserID, uid)
}

func UserIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyUserID).(string); ok {
		return v
	}
	return ""
}
