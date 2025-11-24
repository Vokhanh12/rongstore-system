package db

import (
	"fmt"

	c "server/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitGormPostgresDB(cfg *c.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
