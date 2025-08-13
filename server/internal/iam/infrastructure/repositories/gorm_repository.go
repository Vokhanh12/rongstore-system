// infrastructure/repositories/gorm_repository.go
package repositories

import (
	"context"

	"myapp/internal/iam/domain"

	"gorm.io/gorm"
)

type GormRepository struct {
	rongstore_db *gorm.DB
}

func NewGormRepository(db *gorm.DB) domain.UserRepository {
	return &GormRepository{
		rongstore_db: db,
	}
}

func (r *GormRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	// var user domain.User
	// err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	// if err != nil {
	// 	return nil, err
	// }
	return &domain.User{
		ID: "ID", Email: "EMAIL", Password: "PASSWORD",
	}, nil
}
