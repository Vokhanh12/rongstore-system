package grpc

import (
	"context"
	"time"

	"server/pkg/logger"
	"server/pkg/util/ctxutil"

	"google.golang.org/grpc"
)

func LoggingUnaryInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start).Milliseconds()

		traceID := ctxutil.TraceIDFromContext(ctx)
		ip := peerIP(ctx)
		handlerName := simplifyMethod(info.FullMethod)

		logger.LogAccess(ctx, logger.AccessParams{
			Service:   serviceName,
			Handler:   handlerName,
			Method:    info.FullMethod,
			LatencyMS: duration,
			IP:        ip,
			Status:    statusLabel(err),
			Extra: map[string]interface{}{
				"trace_id": traceID,
			},
		})

		return resp, err
	}
}
