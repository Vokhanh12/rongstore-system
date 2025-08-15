package errors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			appErr := TranslateDomainError(ctx, err)
			// Map HTTP status → gRPC code
			grpcCode := codes.Internal
			if appErr.Status >= 400 && appErr.Status < 500 {
				grpcCode = codes.InvalidArgument
			}
			return nil, status.Errorf(grpcCode, "%s|%s", appErr.Code, appErr.Message)
		}
		return resp, nil
	}
}
