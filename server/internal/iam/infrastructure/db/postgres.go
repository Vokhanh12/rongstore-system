package db

import (
	"context"
	"fmt"
	"time"

	"server/pkg/config"
	"server/pkg/logger"

	domain_errors "server/internal/iam/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitGormPostgresDB(ctx context.Context, cfg *config.Config) *gorm.DB {
	maxRetries := cfg.MaxRetries
	interval := time.Duration(cfg.Interval) * time.Second

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
	)

	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db
		}

		fields := map[string]interface{}{
			"retry":     i + 1,
			"operation": "init.gorm.postgres",
		}

		if i < maxRetries-1 {
			logger.LogInfraDebug(ctx, err, "", fields)
		} else {
			logger.LogBySeverity(ctx, err, fields)
		}

		time.Sleep(interval * time.Duration(1<<i))
	}

	be := domain_errors.POSTGRES_UNAVAILABLE
	panic(fmt.Sprintf(
		"PANIC: [%s][%s] %s | cause: %s | server_action: %s | retryable: %v",
		be.Code,
		be.Key,
		be.Message,
		be.Cause,
		be.ServerAction,
		be.Retryable,
	))
}
