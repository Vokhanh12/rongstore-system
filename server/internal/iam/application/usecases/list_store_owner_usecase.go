package usecases

// import (
// 	"context"
// 	"regexp"
// 	iamv1 "server/api/iam/v1"
// 	"server/internal/iam/domain"
// 	rp "server/internal/iam/domain/repositories"
// 	sv "server/internal/iam/domain/services"
// 	"server/pkg/errors"
// 	"strings"
// )

// type ListStoreOwnerCommand struct {
// 	Email    string
// 	Password string
// }

// type ListStoreOwnerResult struct {
// 	AccessToken  string
// 	RefreshToken string
// }

// func MapListStoreOwnerRequestToCommand(req *iamv1.ListStoreOwnerRequest) ListStoreOwnerCommand {
// 	return ListStoreOwnerCommand{
// 		Email:    req.Email,
// 		Password: req.Password,
// 	}
// }

// func MapListStoreOwnerResultToResponseDTO(result *ListStoreOwnerResult) iamv1.ListStoreOwnerResponse {
// 	return iamv1.ListStoreOwnerResponse{
// 		AccessToken:  result.AccessToken,
// 		RefreshToken: result.RefreshToken,
// 	}
// }

// type ListStoreOwnerUsecase struct {
// 	UserRepo rp.UserRepository
// 	Keycloak sv.Keycloak
// }

// func NewListStoreOwnerUsecase(repo rp.UserRepository, kcl sv.Keycloak) *ListStoreOwnerUsecase {
// 	return &ListStoreOwnerUsecase{
// 		UserRepo: repo,
// 		Keycloak: kcl,
// 	}
// }

// func (u *ListStoreOwnerUsecase) Execute(ctx context.Context, cmd ListStoreOwnerCommand) (*ListStoreOwnerResult, *errors.AppError) {
// 	if cmd.Email == "" {
// 		return nil, errors.New(domain.ListStoreOwner_EMAIL_EMPTY)
// 	}

// 	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
// 	if matched, _ := regexp.MatchString(emailRegex, cmd.Email); !matched {
// 		return nil, errors.New(domain.ListStoreOwner_EMAIL_INVALID)
// 	}

// 	if cmd.Password == "" {
// 		return nil, errors.New(domain.ListStoreOwner_PASSWORD_EMPTY)
// 	}

// 	token, err := u.Keycloak.GetToken(ctx, cmd.Email, cmd.Password)

// 	if err != nil {
// 		if strings.Contains(err.Error(), "invalid_grant") || strings.Contains(err.Error(), "invalid_credentials") {
// 			return nil, errors.New(domain.INVALID_CREDENTIALS)
// 		}
// 		return nil, err
// 	}

// 	return &ListStoreOwnerResult{
// 		AccessToken:  token.AccessToken,
// 		RefreshToken: token.RefreshToken,
// 	}, nil
// }
