package usecases

import (
	"context"
	iamv1 "server/api/iam/v1"
	"server/internal/iam/domain/repositories"
	"server/internal/iam/infrastructure/client"
)

type LoginCommand struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}

func MapLoginRequestToCommand(req *iamv1.LoginRequest) LoginCommand {
	return LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	}
}

func MapLoginResultToResponseDTO(result *LoginResult) iamv1.LoginResponse {
	return iamv1.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}
}

type LoginUserUsecase struct {
	UserRepo repositories.UserRepository
	Keycloak client.KeycloakClient
}

func NewLoginUserUsecase(repo repositories.UserRepository, kcl client.KeycloakClient) *LoginUserUsecase {
	return &LoginUserUsecase{
		UserRepo: repo,
		Keycloak: kcl,
	}
}

func (u *LoginUserUsecase) Execute(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	token, err := u.Keycloak.GetToken(ctx, cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
