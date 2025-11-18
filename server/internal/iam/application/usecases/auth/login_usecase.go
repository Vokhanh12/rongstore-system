package auth

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

// Execute thực hiện đăng nhập
func (u *LoginUsecase) Execute(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	// 0) Validate input
	if cmd.Email == "" {
		return nil, errors.NewBusinessError(domain.LOGIN_EMAIL_EMPTY)
	}

	// RFC5322 basic regex cho email
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, cmd.Email); !matched {
		return nil, errors.NewBusinessError(domain.LOGIN_EMAIL_INVALID)
	}

	if cmd.Password == "" {
		return nil, errors.NewBusinessError(domain.LOGIN_PASSWORD_EMPTY)
	}

	// 1) Gọi Keycloak để lấy token
	token, err := u.Keycloak.GetToken(ctx, cmd.Email, cmd.Password)
	if err != nil {
		// Kiểm tra chuỗi lỗi
		if strings.Contains(err.Error(), "invalid_grant") || strings.Contains(err.Error(), "invalid_credentials") {
			return nil, errors.NewBusinessError(domain.INVALID_CREDENTIALS)
		}
		return nil, err
	}

	// 2) Trả kết quả LoginResult
	return &LoginResult{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
