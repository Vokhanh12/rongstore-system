package queries

import (
	"context"
	en "server/internal/iam/domain/entities"
	rp "server/internal/iam/domain/repositories"
)

type GetUserByEmailQuery struct {
	Email string
}

type GetUserByEmailHandler struct {
	repo rp.UserRepository
}

func NewGetUserByEmailHandler(r rp.UserRepository) *GetUserByEmailHandler {
	return &GetUserByEmailHandler{repo: r}
}

func (h *GetUserByEmailHandler) Handle(ctx context.Context, q GetUserByEmailQuery) (*en.User, error) {
	return h.repo.FindByEmail(ctx, q.Email)
}
