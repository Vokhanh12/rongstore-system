package usecases

import (
	"context"
	"server/internal/iam/application/commands"
	"server/internal/iam/domain"
)

type LoginUserUsecase struct {
	UserRepo domain.UserRepository
	Keycloak domain.Keycloak
}

func NewLoginUserUsecase(repo domain.UserRepository, kcl domain.Keycloak) *LoginUserUsecase {
	return &LoginUserUsecase{
		UserRepo: repo,
		Keycloak: kcl,
	}
}

func (u *LoginUserUsecase) Execute(ctx context.Context, cmd commands.LoginCommand) (*commands.LoginResult, error) {
	token, err := u.Keycloak.GetToken(ctx, cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}

	return &commands.LoginResult{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
