//go:build wireinject
// +build wireinject

package wire

import (
	"myapp/internal/iam/application/usecases"
	"myapp/internal/iam/infrastructure/repositories"
	"myapp/internal/iam/interface/grpc"
	"myapp/pkg/config"

	"github.com/google/wire"
)

func InitializeUserHandler() (*grpc.UserHandler, error) {
	wire.Build(
		config.Load,
		config.NewGormDB,
		repositories.NewGormRepository,
		usecases.NewLoginUserUsecase,
		usecases.NewHandshakeUsecase,
		grpc.NewUserHandler,
	)
	return &grpc.UserHandler{}, nil
}
