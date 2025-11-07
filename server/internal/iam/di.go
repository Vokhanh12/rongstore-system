package wire

import (
	"server/internal/iam/domain/services"
	"server/internal/iam/infrastructure/eventbus"
	iagrpc "server/internal/iam/interface/grpc"
)

type IamDeps struct {
	Handler  *iagrpc.IamHandler
	Store    services.SessionStore
	Keycloak services.Keycloak
	EventBus *eventbus.RabbitMQEventBus
}
