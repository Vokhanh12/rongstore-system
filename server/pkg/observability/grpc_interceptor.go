package observability

import (
	"context"
	"strings"
	"time"

	"server/pkg/errors"
	"server/pkg/metrics"

	"github.com/google/uuid"

	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GrpcTraceUnaryInterceptor extracts "x-trace-id" or "trace-id" from incoming metadata
// and injects it into the context using errors.WithTraceID.
// If no trace id is provided by the client, it generates a UUID and injects it.
func GrpcTraceUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var trace string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get("x-trace-id"); len(vals) > 0 {
				trace = vals[0]
			} else if vals := md.Get("trace-id"); len(vals) > 0 {
				trace = vals[0]
			}
		}
		// generate one if missing
		if trace == "" {
			trace = uuid.New().String()
		}
		ctx = errors.WithTraceID(ctx, trace)
		return handler(ctx, req)
	}
}

// httpStatusToGRPCCode maps HTTP-like statuses to gRPC codes (simple mapping).
func httpStatusToGRPCCode(status int) codes.Code {
	switch {
	case status >= 400 && status < 500:
		return codes.InvalidArgument
	default:
		return codes.Internal
	}
}

// UnaryServerInterceptor wraps handler execution to:
// - track metrics
// - log request/response status
// - translate errors into gRPC status
func UnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// derive handlerName from full method to keep label cardinality low
		fullMethod := info.FullMethod                        // e.g. /package.Service/Method
		handlerName := handlerNameFromFullMethod(fullMethod) // implement helper below

		// metrics: inflight
		metrics.InflightRequests.WithLabelValues(serviceName, handlerName).Inc()
		defer metrics.InflightRequests.WithLabelValues(serviceName, handlerName).Dec()

		resp, err := handler(ctx, req)
		latencyMs := time.Since(start).Milliseconds()

		statusLabel := "OK"
		var appErr *errors.AppError // reuse translated app error if needed
		if err != nil {
			// translate once and reuse
			appErr = errors.TranslateDomainError(ctx, err)
			// if TranslateDomainError returned nil for some reason, make a fallback
			if appErr == nil {
				appErr = &errors.AppError{
					Code:    "UNKNOWN",
					Status:  500,
					Message: "internal error",
					TraceID: errors.TraceIDFromContext(ctx),
				}
			}
			// determine status for metrics/log
			if st, ok := status.FromError(err); ok {
				statusLabel = st.Code().String()
			} else {
				statusLabel = "ERROR"
			}
			// increment error metric by app-level code
			metrics.ErrorsTotal.WithLabelValues(serviceName, handlerName, appErr.Code).Inc()
		}

		// RequestsTotal labels: service, handler, method, status
		metrics.RequestsTotal.WithLabelValues(serviceName, handlerName, fullMethod, statusLabel).Inc()
		// RequestDuration labels: service, handler, method
		metrics.RequestDuration.WithLabelValues(serviceName, handlerName, fullMethod).Observe(float64(latencyMs) / 1000.0)

		// single structured log line per request
		hostname, _ := os.Hostname()
		env := os.Getenv("ENV")
		if env == "" {
			env = "unknown"
		}
		zap.L().Info("grpc_request",
			zap.String("service", serviceName),
			zap.String("handler", handlerName),
			zap.String("method", fullMethod),
			zap.String("trace_id", errors.TraceIDFromContext(ctx)),
			zap.String("status", statusLabel),
			zap.Int64("latency_ms", latencyMs),
			zap.String("instance", hostname),
			zap.String("env", env),
		)

		// If no error, return response
		if err == nil {
			return resp, nil
		}

		// convert domain error -> grpc status using the already-translated appErr
		grpcCode := httpStatusToGRPCCode(appErr.Status)
		st := status.New(grpcCode, strings.Join([]string{appErr.Code, appErr.Message}, "|"))
		return nil, st.Err()
	}
}

// handlerNameFromFullMethod extracts a short handler name from info.FullMethod.
// e.g. "/package.Service/Method" -> "Service/Method"
func handlerNameFromFullMethod(fullMethod string) string {
	if fullMethod == "" {
		return "unknown"
	}
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	parts := strings.SplitN(fullMethod, "/", 2)
	if len(parts) == 2 {
		servicePart := parts[0]
		if idx := strings.LastIndex(servicePart, "."); idx != -1 && idx+1 < len(servicePart) {
			servicePart = servicePart[idx+1:]
		}
		// return "Service/Method" to keep label low-cardinality
		return servicePart + "/" + parts[1]
	}
	if idx := strings.LastIndex(fullMethod, "."); idx != -1 && idx+1 < len(fullMethod) {
		return fullMethod[idx+1:]
	}
	return fullMethod
}
