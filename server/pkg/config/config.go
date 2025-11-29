package config

import (
	"log"
	"os"
	"strconv"

	"server/internal/iam/domain"
	"server/pkg/errors"

	"github.com/joho/godotenv"
)

type Config struct {

	// DB
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	// Redis
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int
	RedisTTL      int

	// Keycloak
	KeycloakURL                     string
	KeycloakServerHealth            string
	KeycloakRealm                   string
	KeycloakClientID                string
	KeycloakSecret                  string
	KeycloakRedirect                string
	KeycloakScope                   string
	KeycloakGrantUmaTicketType      string
	KeycloakAudience                string
	KeycloakResponsePermissionsMode string

	// Title-gl
	TitleGlHost string
	TitleGlPort int

	// ORSM
	ORSMHost string
	ORSMPort string

	// RabitMQ
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUser     string
	RabbitMQPassword string

	MaxRetries int
	Interval   int

	IAMDefaultErrors map[string]errors.BusinessError
}

func Load() *Config {
	// Load .env nếu có
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, reading env variables from system")
	}

	port, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid POSTGRES_PORT: %v", err)
	}

	redisPort, err := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	if err != nil {
		log.Fatalf("Invalid REDIS_PORT: %v", err)
	}
	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		log.Fatalf("Invalid REDIS_DB: %v", err)
	}

	redisTTL, err := strconv.Atoi(getEnv("REDIS_TTL", "900"))
	if err != nil {
		log.Fatalf("Invalid REDIS_TTL: %v", err)
	}

	maxRetries, err := strconv.Atoi(getEnv("MAX_RETRIES", "3"))
	if err != nil {
		log.Fatalf("Invalid MAX_RETRIES: %v", err)
	}

	interval, err := strconv.Atoi(getEnv("INTERVAL", "1"))
	if err != nil {
		log.Fatalf("Invalid INTERVAL: %v", err)
	}

	return &Config{
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     port,
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       getEnv("POSTGRES_DB", "postgres"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     redisPort,
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
		RedisTTL:      redisTTL,

		KeycloakURL:                     getEnv("KEYCLOAK_SERVER_URL", ""),
		KeycloakServerHealth:            getEnv("KEYCLOAK_SERVER_HEALTH", ""),
		KeycloakRealm:                   getEnv("KEYCLOAK_REALM", ""),
		KeycloakClientID:                getEnv("KEYCLOAK_CLIENT_ID", ""),
		KeycloakSecret:                  getEnv("KEYCLOAK_CLIENT_SECRET", ""),
		KeycloakRedirect:                getEnv("KEYCLOAK_REDIRECT_URI", ""),
		KeycloakScope:                   getEnv("KEYCLOAK_SCOPE", ""),
		KeycloakGrantUmaTicketType:      getEnv("KEYCLOAK_GRANT_UMA_TICKET_TYPE", ""),
		KeycloakAudience:                getEnv("KEYCLOAK_AUDIENCE", ""),
		KeycloakResponsePermissionsMode: getEnv("KEYCLOAK_RESPONSE_PERMISSIONS_MODE", ""),

		RabbitMQHost:     getEnv("RABBITMQ_HOST", ""),
		RabbitMQPort:     getEnv("RABBITMQ_PORT", ""),
		RabbitMQUser:     getEnv("RABBITMQ_USER", ""),
		RabbitMQPassword: getEnv("RABBITMQ_PASSWORD", ""),

		MaxRetries: maxRetries,
		Interval:   interval,

		IAMDefaultErrors: domain.ErrorByCode,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
