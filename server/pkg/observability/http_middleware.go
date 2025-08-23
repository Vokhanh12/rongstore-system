package observability

import (
	"net/http"
	"time"

	"server/pkg/errors"

	"server/internal/iam/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func HTTPMiddleware(next http.Handler, store domain.SessionStore, serviceName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()

		trace := r.Header.Get("X-Trace-Id")
		if trace == "" {
			trace = r.Header.Get("Trace-Id")
		}
		if trace == "" {
			trace = uuid.NewString()
		}
		ctx = errors.WithTraceID(ctx, trace)

		// session id
		sid := r.Header.Get("X-Session-Id")
		if sid == "" {
			sid = r.Header.Get("Session-Id")
		}
		if sid != "" {
			ctx = errors.WithSessionID(ctx, sid)
			if store != nil {
				if s, _ := store.GetSession(ctx, sid); s != nil {
					ctx = errors.WithUserID(ctx, s.UserID)
				}
			}
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		latency := time.Since(start).Milliseconds()
		fields := []zap.Field{
			zap.String("service", serviceName),
			zap.String("handler", r.URL.Path),
			zap.String("method", r.Method),
			zap.String("trace_id", errors.TraceIDFromContext(ctx)),
			zap.Int64("latency_ms", latency),
			zap.String("instance", hostnameCache),
			//zap.String("env", func() string { e := os.Getenv("ENV"); if e==\"\" { return \"unknown\" }; return e }()),
		}
		if sid := errors.SessionIDFromContext(ctx); sid != "" {
			fields = append(fields, zap.String("session_id", sid))
		}
		if uid := errors.UserIDFromContext(ctx); uid != "" {
			fields = append(fields, zap.String("user_id", uid))
		}

		zap.L().With(fields...).Info("http_request")
	})
}
