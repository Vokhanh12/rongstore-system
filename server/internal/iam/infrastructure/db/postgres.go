package db

import (
	"context"
	"fmt"
	"time"

	"server/pkg/config"
	"server/pkg/errors"
	"server/pkg/logger"
	"server/pkg/util/infahelper"

	domain_errors "server/internal/iam/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitGormPostgresDB(
	ctx context.Context,
	cfg *config.Config,
) *gorm.DB {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
	)

	db, err := infahelper.Retry(
		cfg.MaxRetries,
		time.Duration(cfg.Interval)*time.Second,
		func() (*gorm.DB, *errors.AppError) {
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				return nil, errors.New(
					domain_errors.POSTGRES_UNAVAILABLE,
					errors.SetError(err),
				)
			}
			return db, nil
		},
	)

	if err != nil {
		logger.LogBySeverity(
			ctx,
			"Init.Postgres",
			err,
		)
		return nil
	}

	return db
}
