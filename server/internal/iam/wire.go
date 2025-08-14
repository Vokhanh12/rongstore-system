//go:build wireinject
// +build wireinject

package wire

import (
	"server/internal/iam/application/usecases"
	"server/internal/iam/infrastructure/repositories"
	"server/internal/iam/interface/grpc"
	"server/pkg/config"

	"github.com/google/wire"
)

func InitializeIamHandler() (*grpc.IamHandler, error) {
	wire.Build(
		config.Load,
		config.NewGormDB,
		repositories.NewGormRepository,
		usecases.NewLoginUserUsecase,
		usecases.NewHandshakeUsecase,
		grpc.NewIamHandler,
	)
	return &grpc.IamHandler{}, nil
}
