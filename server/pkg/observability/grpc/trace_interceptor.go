package grpc

import (
	"context"

	"server/pkg/util/ctxutil"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TraceUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		traceID := extractTraceID(ctx)
		ctx = ctxutil.WithTraceID(ctx, traceID)

		sessionID := extractSessionID(ctx)
		ctx = ctxutil.WithSessionID(ctx, sessionID)

		requestID := extractRequestID(ctx)
		ctx = ctxutil.WithRequestID(ctx, requestID)

		return handler(ctx, req)
	}
}

func extractTraceID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for _, key := range []string{"x-trace-id", "trace-id"} {
			if vals := md.Get(key); len(vals) > 0 {
				return vals[0]
			}
		}
	}
	return uuid.NewString()
}

func extractRequestID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for _, key := range []string{"x-request-id", "request-id"} {
			if vals := md.Get(key); len(vals) > 0 {
				return vals[0]
			}
		}
	}
	return uuid.NewString()
}

func extractSessionID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for _, key := range []string{"x-session-id", "session-id"} {
			if vals := md.Get(key); len(vals) > 0 {
				return vals[0]
			}
		}
	}
	return uuid.NewString()
}
