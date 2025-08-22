//go:build wireinject
// +build wireinject

package wire

import (
	"server/internal/iam/application/usecases"
	"server/internal/iam/infrastructure/cache"
	"server/internal/iam/infrastructure/repositories"
	"server/internal/iam/interface/grpc"
	"server/pkg/config"
	"server/pkg/observability"

	"github.com/google/wire"
)

func InitializeIamHandler() (*grpc.IamHandler, observability.SessionStore, error) {
	wire.Build(
		config.Load,
		config.NewRedisClient,
		config.NewGormDB,
		repositories.NewGormRepository,
		cache.RedisTTLFromConfig,
		cache.NewRedisSessionStoreProvider,        // -> *cache.RedisSessionStore
		cache.NewObservabilitySessionStoreAdapter, // -> observability.SessionStore
		usecases.NewLoginUserUsecase,
		usecases.NewHandshakeUsecase,
		grpc.NewIamHandler,
	)
	return &grpc.IamHandler{}, nil, nil
}
