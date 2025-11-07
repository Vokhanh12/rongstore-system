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
