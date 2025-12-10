package trace

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey string

const traceIDKey = ctxKey("trace_id")
const TraceHeader = "X-Request-Id"

func NewTraceID() string {
	return uuid.NewString()
}

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey, id)
}

func TraceIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(traceIDKey)
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(TraceHeader)
		if id == "" {
			id = NewTraceID()
		}
		w.Header().Set(TraceHeader, id)
		ctx := WithTraceID(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
