package observability

import (
	"context"
	"time"

	"server/pkg/errors"
	"server/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryLoggingInterceptor logs each unary RPC with trace id, method, latency and status.
// Keep it lightweight: avoid expensive serialization. req/resp are logged with zap.Any (beware PII/secrets).
func UnaryLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()

		// extract trace id if any
		traceID := errors.TraceIDFromContext(ctx)
		// also attempt to extract metadata request id if provided
		var requestID string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if v := md.Get("x-request-id"); len(v) > 0 {
				requestID = v[0]
			}
		}

		// Call actual handler
		resp, err = handler(ctx, req)
		elapsed := time.Since(start)

		// status
		statusCode := "OK"
		if err != nil {
			if st, _ := status.FromError(err); st != nil {
				statusCode = st.Code().String()
			} else {
				statusCode = "ERROR"
			}
		}

		// prepare zap fields
		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("trace_id", traceID),
			zap.String("status", statusCode),
			zap.Duration("latency_ms", elapsed),
		}
		if requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		// Optionally include request/response for debug level - be careful with PII
		// Use logger.L.Check to avoid computing zap.Any if not enabled.
		if ce := logger.L.Check(zap.DebugLevel, "rpc"); ce != nil {
			// redact payload before logging (basic)
			ce.Write(append(fields, zap.Any("req", redactPayload(req)), zap.Any("resp", redactPayload(resp)))...)
		} else {
			// Info level: no raw payloads
			if err != nil {
				logger.L.Warn("rpc error", fields...)
			} else {
				logger.L.Info("rpc", fields...)
			}
		}

		return resp, err
	}
}

// StreamLoggingInterceptor logs streaming RPCs (open/close) with basic metadata.
func StreamLoggingInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		start := time.Now()
		traceID := errors.TraceIDFromContext(ctx)

		// call handler
		err := handler(srv, ss)
		elapsed := time.Since(start)

		statusCode := "OK"
		if err != nil {
			if st, _ := status.FromError(err); st != nil {
				statusCode = st.Code().String()
			} else {
				statusCode = "ERROR"
			}
		}

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("trace_id", traceID),
			zap.String("status", statusCode),
			zap.Duration("latency_ms", elapsed),
			zap.Bool("streaming", info.IsServerStream || info.IsClientStream),
		}

		if err != nil {
			logger.L.Warn("grpc stream error", fields...)
		} else {
			logger.L.Info("grpc stream", fields...)
		}
		return err
	}
}

// -----------------------------------------------------------------------------
// Handler-level helper: log contextual business events with trace id
// Usage inside handler/usecase: observability.LogWithTrace(ctx).Info("user created", zap.String("user", id))
// -----------------------------------------------------------------------------

// LogWithTrace returns a sugared logger prefilled with trace_id so handler can log business events.
func LogWithTrace(ctx context.Context) *zap.Logger {
	traceID := errors.TraceIDFromContext(ctx)
	// create a child logger with trace_id field
	return logger.L.With(zap.String("trace_id", traceID))
}

// redactPayload attempts to mask sensitive fields in common request/response shapes.
// It's intentionally conservative: if payload is map[string]interface{} or struct it will try
// to mask keys like "password", "token", "access_token", "refresh_token".
// You should replace/extend with schema-specific redaction logic.
func redactPayload(v interface{}) interface{} {
	// very small, non-recursive heuristic to avoid leaking secrets in debug logs.
	// If you have proto messages, you can implement a more correct mechanism (marshal -> map -> redact).
	switch t := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(t))
		for k, val := range t {
			lk := lower(k)
			if isSensitiveKey(lk) {
				out[k] = "<REDACTED>"
			} else {
				out[k] = val
			}
		}
		return out
	default:
		// for other types (proto messages, structs) we avoid deep introspection by default.
		// Returning the original value - but we only call this in Debug level to limit exposure.
		return v
	}
}

func lower(s string) string {
	// small helper
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] = b[i] + 32
		}
	}
	return string(b)
}

func isSensitiveKey(k string) bool {
	switch k {
	case "password", "passwd", "token", "access_token", "refresh_token", "secret", "ssn":
		return true
	default:
		return false
	}
}
