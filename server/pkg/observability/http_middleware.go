package observability

import (
	"net/http"
	"time"

	"server/pkg/errors"
	"server/pkg/metrics"

	"go.uber.org/zap"
)

func HTTPMiddleware(serviceName, handlerName, method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		metrics.InflightRequests.WithLabelValues(serviceName, handlerName).Inc()
		defer metrics.InflightRequests.WithLabelValues(serviceName, handlerName).Dec()

		// propagate trace id if present; errors.WithTraceID can be used in handlers
		trace := r.Header.Get("X-Trace-Id")
		if trace != "" {
			r = r.WithContext(errors.WithTraceID(r.Context(), trace))
		}

		// wrapper to capture status code
		ww := &responseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)
		duration := time.Since(start).Seconds()

		codeLabel := "OK"
		if ww.status >= 400 {
			codeLabel = http.StatusText(ww.status)
		}

		metrics.RequestsTotal.WithLabelValues(serviceName, handlerName, method, codeLabel).Inc()
		metrics.RequestDuration.WithLabelValues(serviceName, handlerName, method).Observe(duration)

		// optional structured log
		zap.L().Info("http_request",
			zap.String("service", serviceName),
			zap.String("handler", handlerName),
			zap.String("method", method),
			zap.Int("status", ww.status),
			zap.Float64("duration_s", duration),
			zap.String("trace_id", errors.TraceIDFromContext(r.Context())),
		)
	})
}

// simple wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
