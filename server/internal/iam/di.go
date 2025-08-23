package wire

import (
	"server/internal/iam/domain"
	iagrpc "server/internal/iam/interface/grpc"
)

type IamDeps struct {
	Handler *iagrpc.IamHandler
	Store   domain.SessionStore
}
