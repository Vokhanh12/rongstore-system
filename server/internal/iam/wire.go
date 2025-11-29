//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"server/internal/iam/application/usecases"
	"server/internal/iam/infrastructure/cache"
	"server/internal/iam/infrastructure/client"
	"server/internal/iam/infrastructure/db"
	domain_errors "server/internal/iam/infrastructure/errors"
	"server/internal/iam/infrastructure/eventbus"
	"server/internal/iam/infrastructure/repositories"
	"server/internal/iam/interface/grpc"
	"server/pkg/config"

	"github.com/google/wire"
)

func InitializeIamHandler(ctx context.Context) IamDeps {
	wire.Build(
		config.Load,
		domain_errors.InitBusinessError,
		cache.InitRedisSessionStore,
		client.InitKeycloakClient,
		eventbus.InitRabbitMQEventBus,
		db.InitGormPostgresDB,
		repositories.NewGormUserRepository,
		usecases.NewLoginUsecase,
		usecases.NewHandshakeUsecase,
		grpc.NewIamHandler,
		wire.Struct(new(IamDeps), "Handler", "RedisSessionStore", "Keycloak", "EventBus", "BusinessError"),
	)
	return IamDeps{}
}
