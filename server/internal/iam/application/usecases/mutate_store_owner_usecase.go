package usecases

import (
	"context"
	commonv1 "server/api/common/v1"
	resourcesv1 "server/api/rongstore/v1/resources"
	rp "server/internal/iam/domain/repositories"
	sv "server/internal/iam/domain/services"
	"server/pkg/errors"
)

type MutationType int

const (
	MutationCreate MutationType = iota
	MutationUpdate
	MutationDelete
)

type StoreOwner struct {
	OrderID     string
	LineNumber  int32
	ProductCode string
	ProductName string
}

type StoreOwnerMutateItem struct {
	Type  MutationType
	Value StoreOwner
}

type StoreOwnerMutateCommand struct {
	Items []StoreOwnerMutateItem
}

type MutateResultItem struct {
	Success bool
	Name    string
	Code    string
	Error   errors.AppError
	Details map[string]string
}

type StoreOwnerMutateResult struct {
	Items []MutateResultItem
}

func mapStoreOwner(p *resourcesv1.StoreOwner) StoreOwner {
	return StoreOwner{
		OrderID:     p.OrderId,
		LineNumber:  p.LineNumber,
		ProductCode: p.ProductCode,
		ProductName: p.ProductName,
	}
}

func MapStoreOwnerMutateRequestToCommand(
	req *resourcesv1.StoreOwnerMutateRequest,
) StoreOwnerMutateCommand {

	items := make([]StoreOwnerMutateItem, 0)

	for _, m := range req.Data.Mutations {
		switch v := m.Mutation.(type) {

		case *resourcesv1.StoreOwnerMutateRequest_StoreOwnerMutateOneOf_Create:
			items = append(items, StoreOwnerMutateItem{
				Type:  MutationCreate,
				Value: mapStoreOwner(v.Create),
			})

		case *resourcesv1.StoreOwnerMutateRequest_StoreOwnerMutateOneOf_Update:
			items = append(items, StoreOwnerMutateItem{
				Type:  MutationUpdate,
				Value: mapStoreOwner(v.Update),
			})

		case *resourcesv1.StoreOwnerMutateRequest_StoreOwnerMutateOneOf_Delete:
			items = append(items, StoreOwnerMutateItem{
				Type:  MutationDelete,
				Value: mapStoreOwner(v.Delete),
			})
		}
	}

	return StoreOwnerMutateCommand{Items: items}
}

func MapStoreOwnerMutateResultToResponseDTO(
	result *StoreOwnerMutateResult,
) *commonv1.MutateResponse {

	items := make([]*commonv1.MutateResult, 0, len(result.Items))

	for _, item := range result.Items {
		items = append(items, &commonv1.MutateResult{
			Success: item.Success,
			Name:    item.Name,
			Error:   MapAppErrorToProto(item.Error),
			Details: item.Details,
		})
	}

	return &commonv1.MutateResponse{
		Data: &commonv1.MutateResponse_MutateResponseData{
			MutateResult: items,
		},
	}
}

type StoreOwnerMutateUsecase struct {
	UserRepo rp.UserRepository
	Keycloak sv.Keycloak
}

func NewMutateStoreOwnerUsecase(repo rp.UserRepository, kcl sv.Keycloak) *StoreOwnerMutateUsecase {
	return &StoreOwnerMutateUsecase{
		UserRepo: repo,
		Keycloak: kcl,
	}
}

func (u *StoreOwnerMutateUsecase) Execute(ctx context.Context, cmd StoreOwnerMutateCommand) (*StoreOwnerMutateResult, *errors.AppError) {

	results := make([]MutateResultItem, 0)

	for _, item := range cmd.Items {
		var (
			err error
			id  string
		)

		switch item.Type {

		case MutationCreate:
			id, err = u.Repo.Create(ctx, item.Value)

		case MutationUpdate:
			err = u.Repo.Update(ctx, item.Value)
			id = item.Value.OrderID

		case MutationDelete:
			err = u.Repo.Delete(ctx, item.Value)
			id = item.Value.OrderID
		}

		results = append(results, MutateResultItem{
			Name: id,
			Code: errors.ParseError(err).Error(),
		})
	}

	return &StoreOwnerMutateResult{
		Items: results,
	}, nil
}
