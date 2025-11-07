//go:build wireinject
// +build wireinject

package wire

import (
	"server/internal/iam/application/usecases/auth"
	"server/internal/iam/infrastructure/cache"
	"server/internal/iam/infrastructure/client"
	"server/internal/iam/infrastructure/eventbus"
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
		eventbus.NewEventBusFromConfig,
		repositories.NewGormUserRepository,
		cache.RedisTTLFromConfig,
		cache.NewRedisSessionStore,
		auth.NewLoginUsecase,
		auth.NewHandshakeUsecase,
		grpc.NewIamHandler,
		wire.Struct(new(IamDeps), "Handler", "Store", "Keycloak", "EventBus"),
	)
	return IamDeps{}, nil
}
