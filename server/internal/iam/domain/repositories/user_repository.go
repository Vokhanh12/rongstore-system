package repositories

import (
	"context"
	"server/internal/iam/domain/entities"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
}
