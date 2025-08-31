package observability

import (
	"context"
	"time"

	"server/pkg/logger"
	"server/pkg/util/ctxutil"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()

		traceID := ctxutil.GetIDFromContext(ctx)
		var requestID string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if v := md.Get("x-request-id"); len(v) > 0 {
				requestID = v[0]
			}
		}

		resp, err = handler(ctx, req)
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
		}
		if requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		if ce := logger.L.Check(zap.DebugLevel, "rpc"); ce != nil {
			ce.Write(append(fields, zap.Any("req", redactPayload(req)), zap.Any("resp", redactPayload(resp)))...)
		} else {
			if err != nil {
				logger.L.Warn("rpc error", fields...)
			} else {
				logger.L.Info("rpc", fields...)
			}
		}

		return resp, err
	}
}

func StreamLoggingInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		start := time.Now()
		traceID := ctxutil.GetIDFromContext(ctx)

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

func LogWithTrace(ctx context.Context) *zap.Logger {
	traceID := ctxutil.GetIDFromContext(ctx)
	return logger.L.With(zap.String("trace_id", traceID))
}

func redactPayload(v interface{}) interface{} {
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
		return v
	}
}

func lower(s string) string {
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
