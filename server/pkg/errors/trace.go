package errors

import "context"

type traceKey struct{}

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceKey{}, id)
}

func TraceIDFromContext(ctx context.Context) string {
	if v := ctx.Value(traceKey{}); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
