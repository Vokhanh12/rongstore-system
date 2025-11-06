package repositories

import (
	"context"

	en "server/internal/iam/domain/entities"
	rp "server/internal/iam/domain/repositories"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	rongstore_db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) rp.UserRepository {
	return &GormUserRepository{
		rongstore_db: db,
	}
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*en.User, error) {
	// var user domain.User
	// err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	// if err != nil {
	// 	return nil, err
	// }
	return &en.User{
		ID: "ID", Email: "EMAIL", Password: "PASSWORD",
	}, nil
}
