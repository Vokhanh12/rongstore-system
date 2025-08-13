package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	// DB config
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	// Keycloak config
	KeycloakURL      string
	KeycloakRealm    string
	KeycloakClientID string
	KeycloakSecret   string
	KeycloakRedirect string
	KeycloakScope    string
}

func Load() *Config {
	// Load .env nếu có
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println(".env file not found, reading env variables from system")
	}

	port, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid POSTGRES_PORT: %v", err)
	}

	return &Config{
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     port,
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       getEnv("POSTGRES_DB", "postgres"),

		KeycloakURL:      getEnv("KEYCLOAK_SERVER_URL", ""),
		KeycloakRealm:    getEnv("KEYCLOAK_REALM", ""),
		KeycloakClientID: getEnv("KEYCLOAK_CLIENT_ID", ""),
		KeycloakSecret:   getEnv("KEYCLOAK_CLIENT_SECRET", ""),
		KeycloakRedirect: getEnv("KEYCLOAK_REDIRECT_URI", ""),
		KeycloakScope:    getEnv("KEYCLOAK_SCOPE", ""),
	}
}

func NewGormDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
