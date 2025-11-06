package wire

import (
	"server/internal/iam/domain/services"
	iagrpc "server/internal/iam/interface/grpc"
)

type IamDeps struct {
	Handler  *iagrpc.IamHandler
	Store    services.SessionStore
	Keycloak services.Keycloak
}
