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

// type MutateStoreOwnerCommand struct {
// 	Email    string
// 	Password string
// }

// type MutateStoreOwnerResult struct {
// 	AccessToken  string
// 	RefreshToken string
// }

// func MapMutateStoreOwnerRequestToCommand(req *iamv1.MutateStoreOwnerRequest) MutateStoreOwnerCommand {
// 	return MutateStoreOwnerCommand{
// 		Email:    req.Email,
// 		Password: req.Password,
// 	}
// }

// func MapMutateStoreOwnerResultToResponseDTO(result *MutateStoreOwnerResult) iamv1.MutateStoreOwnerResponse {
// 	return iamv1.MutateStoreOwnerResponse{
// 		AccessToken:  result.AccessToken,
// 		RefreshToken: result.RefreshToken,
// 	}
// }

// type MutateStoreOwnerUsecase struct {
// 	UserRepo rp.UserRepository
// 	Keycloak sv.Keycloak
// }

// func NewMutateStoreOwnerUsecase(repo rp.UserRepository, kcl sv.Keycloak) *MutateStoreOwnerUsecase {
// 	return &MutateStoreOwnerUsecase{
// 		UserRepo: repo,
// 		Keycloak: kcl,
// 	}
// }

// func (u *MutateStoreOwnerUsecase) Execute(ctx context.Context, cmd MutateStoreOwnerCommand) (*MutateStoreOwnerResult, *errors.AppError) {
// 	if cmd.Email == "" {
// 		return nil, errors.New(domain.MutateStoreOwner_EMAIL_EMPTY)
// 	}

// 	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
// 	if matched, _ := regexp.MatchString(emailRegex, cmd.Email); !matched {
// 		return nil, errors.New(domain.MutateStoreOwner_EMAIL_INVALID)
// 	}

// 	if cmd.Password == "" {
// 		return nil, errors.New(domain.MutateStoreOwner_PASSWORD_EMPTY)
// 	}

// 	token, err := u.Keycloak.GetToken(ctx, cmd.Email, cmd.Password)

// 	if err != nil {
// 		if strings.Contains(err.Error(), "invalid_grant") || strings.Contains(err.Error(), "invalid_credentials") {
// 			return nil, errors.New(domain.INVALID_CREDENTIALS)
// 		}
// 		return nil, err
// 	}

// 	return &MutateStoreOwnerResult{
// 		AccessToken:  token.AccessToken,
// 		RefreshToken: token.RefreshToken,
// 	}, nil
// }
