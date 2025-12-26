package grpc

import (
	"context"

	"server/internal/iam/domain/services"
	"server/pkg/auth"
	"server/pkg/util/ctxutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthUnaryInterceptor(store services.IRedisCache, rules []auth.GrpcRule) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		if err := auth.ValidateWithMetadata(md, info.FullMethod, rules); err != nil {
			return nil, err
		}

		if sid := ctxutil.SessionIDFromContext(ctx); sid != "" && store != nil {
			if se, err := store.GetSession(ctx, sid); err == nil && se != nil {
				ctx = ctxutil.WithUserID(ctx, se.UserID)
			}
		}
		return handler(ctx, req)
	}
}
