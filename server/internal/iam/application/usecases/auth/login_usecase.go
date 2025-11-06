package auth

import (
	"context"
	iamv1 "server/api/iam/v1"
	rp "server/internal/iam/domain/repositories"
	sv "server/internal/iam/domain/services"
)

// --- Command & Result (trước đây ở package commands) ---
type LoginCommand struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}

// --- Mapper (trước đây ở package mappers) ---
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

// --- Usecase (trước đây ở package usecases) ---
type LoginUsecase struct {
	UserRepo rp.UserRepository
	Keycloak sv.Keycloak
}

func NewLoginUsecase(repo rp.UserRepository, kcl sv.Keycloak) *LoginUsecase {
	return &LoginUsecase{
		UserRepo: repo,
		Keycloak: kcl,
	}
}

func (u *LoginUsecase) Execute(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	token, err := u.Keycloak.GetToken(ctx, cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
