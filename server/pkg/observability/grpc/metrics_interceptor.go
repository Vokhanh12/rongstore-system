package grpc

import (
	"context"
	"time"

	"server/pkg/metrics"

	"google.golang.org/grpc"
)

func MetricsUnaryInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		method := simplifyMethod(info.FullMethod)

		metrics.InflightRequests.WithLabelValues(serviceName, method).Inc()
		defer metrics.InflightRequests.WithLabelValues(serviceName, method).Dec()

		resp, err := handler(ctx, req)

		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues(serviceName, method, info.FullMethod).Observe(duration)

		return resp, err
	}
}
