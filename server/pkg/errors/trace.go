package errors

import "context"

// traceKey is an unexported type to avoid collisions in context keys.
type traceKey struct{}

// WithTraceID returns a new context with the trace id attached.
func WithTraceID(ctx context.Context, id string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, traceKey{}, id)
}

// TraceIDFromContext extracts trace id from context. Returns empty string if not present.
func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(traceKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
