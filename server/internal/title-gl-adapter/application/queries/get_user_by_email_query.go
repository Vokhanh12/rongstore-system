package queries

import (
	"context"
	"server/internal/iam/domain/entities"
	"server/internal/iam/domain/repositories"
)

type GetUserByEmailQuery struct {
	Email string
}

type GetUserByEmailHandler struct {
	repo repositories.UserRepository
}

func NewGetUserByEmailHandler(r repositories.UserRepository) *GetUserByEmailHandler {
	return &GetUserByEmailHandler{repo: r}
}

func (h *GetUserByEmailHandler) Handle(ctx context.Context, q GetUserByEmailQuery) (*entities.User, error) {
	return h.repo.FindByEmail(ctx, q.Email)
}
