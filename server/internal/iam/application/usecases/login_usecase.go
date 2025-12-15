package usecases

import (
	"context"
	"regexp"
	iamv1 "server/api/iam/v1"
	"server/internal/iam/domain"
	rp "server/internal/iam/domain/repositories"
	sv "server/internal/iam/domain/services"
	"server/pkg/errors"
	"strings"
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

func (u *LoginUsecase) Execute(ctx context.Context, cmd LoginCommand) (*LoginResult, *errors.AppError) {
	if cmd.Email == "" {
		return nil, errors.New(domain.LOGIN_EMAIL_EMPTY)
	}

	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, cmd.Email); !matched {
		return nil, errors.New(domain.LOGIN_EMAIL_INVALID)
	}

	if cmd.Password == "" {
		return nil, errors.New(domain.LOGIN_PASSWORD_EMPTY)
	}

	token, err := u.Keycloak.GetToken(ctx, cmd.Email, cmd.Password)

	if err != nil {
		if strings.Contains(err.Error(), "invalid_grant") || strings.Contains(err.Error(), "invalid_credentials") {
			return nil, errors.New(domain.INVALID_CREDENTIALS)
		}
		return nil, err
	}

	return &LoginResult{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
