//go:build wireinject
// +build wireinject

package wire

import (
	"server/internal/iam/application/usecases"
	"server/internal/iam/infrastructure/cache"
	"server/internal/iam/infrastructure/client"
	"server/internal/iam/infrastructure/repositories"
	"server/internal/iam/interface/grpc"
	"server/pkg/config"

	"github.com/google/wire"
)

func InitializeIamHandler() (IamDeps, error) {
	wire.Build(
		config.Load,
		config.NewRedisClient,
		config.NewGormDB,
		client.NewKeycloakClient,
		repositories.NewGormRepository,
		cache.RedisTTLFromConfig,
		cache.NewRedisSessionStore,
		usecases.NewLoginUserUsecase,
		usecases.NewHandshakeUsecase,
		grpc.NewIamHandler,
		wire.Struct(new(IamDeps), "Handler", "Store", "Keycloak"), // thêm field Keycloak
	)
	return IamDeps{}, nil
}
