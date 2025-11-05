package wire

import (
	"server/internal/title-gl-adapter/domain"
	iagrpc "server/internal/title-gl-adapter/interface/grpc"
)

type IamDeps struct {
	Handler *iagrpc.IamHandler
	title   domain.TitleGl
}
