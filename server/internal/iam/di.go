package wire

import (
	"server/internal/iam/domain/services"
	"server/internal/iam/infrastructure/eventbus"
	iagrpc "server/internal/iam/interface/grpc"
)

type IamDeps struct {
	Handler           *iagrpc.IamHandler
	RedisSessionStore services.RedisSessionStore
	Keycloak          services.Keycloak
	BusinessError     services.BusinessError
	EventBus          *eventbus.RabbitMQEventBus
}
