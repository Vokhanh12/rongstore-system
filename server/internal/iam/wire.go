//go:build wireinject
// +build wireinject

package wire

import (
	"server/internal/iam/application/usecases/auth"
	"server/internal/iam/infrastructure/cache"
	"server/internal/iam/infrastructure/client"
	"server/internal/iam/infrastructure/db"
	"server/internal/iam/infrastructure/eventbus"
	"server/internal/iam/infrastructure/repositories"
	"server/internal/iam/interface/grpc"
	"server/pkg/config"

	"github.com/google/wire"
)

func InitializeIamHandler() (IamDeps, error) {
	wire.Build(
		config.Load,
		cache.InitRedisSessionStore,
		client.InitKeycloakClient,
		eventbus.InitRabbitMQEventBus,
		db.InitGormPostgresDB,
		repositories.NewGormUserRepository,
		auth.NewLoginUsecase,
		auth.NewHandshakeUsecase,
		grpc.NewIamHandler,
		wire.Struct(new(IamDeps), "Handler", "RedisSessionStore", "Keycloak", "EventBus"),
	)
	return IamDeps{}, nil
}
