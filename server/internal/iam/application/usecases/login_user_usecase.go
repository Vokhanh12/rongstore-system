package usecases

import (
	"context"
	"server/internal/iam/application/commands"
	"server/internal/iam/domain"
)

type LoginUserUsecase struct {
	UserRepo domain.UserRepository
}

func NewLoginUserUsecase(repo domain.UserRepository) *LoginUserUsecase {
	return &LoginUserUsecase{
		UserRepo: repo,
	}
}

func (u *LoginUserUsecase) Execute(ctx context.Context, cmd commands.LoginCommand) (*commands.LoginResult, error) {

	_, err := u.UserRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}

	return &commands.LoginResult{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
	}, nil
}
