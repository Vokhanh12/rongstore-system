package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	KeycloakURL          string
	KeycloakServerHealth string
	KeycloakRealm        string
	KeycloakClientID     string
	KeycloakSecret       string
	KeycloakRedirect     string
	KeycloakScope        string

	// Title-gl
	TitleGlHost string
	TitleGlPort int

	// ORSM
	ORSMHost string
	ORSMPort string
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
	redisTTL, err := strconv.Atoi(getEnv("REDIS_TTL", "900")) // mặc định 15 phút
	if err != nil {
		log.Fatalf("Invalid REDIS_TTL: %v", err)
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

		KeycloakURL:          getEnv("KEYCLOAK_SERVER_URL", ""),
		KeycloakServerHealth: getEnv("KEYCLOAK_SERVER_HEALTH", ""),
		KeycloakRealm:        getEnv("KEYCLOAK_REALM", ""),
		KeycloakClientID:     getEnv("KEYCLOAK_CLIENT_ID", ""),
		KeycloakSecret:       getEnv("KEYCLOAK_CLIENT_SECRET", ""),
		KeycloakRedirect:     getEnv("KEYCLOAK_REDIRECT_URI", ""),
		KeycloakScope:        getEnv("KEYCLOAK_SCOPE", ""),
	}
}

func NewGormDB(cfg *Config) (*gorm.DB, error) {
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

// NewRedisClient khởi tạo Redis client từ cấu hình
func NewRedisClient(cfg *Config) *redis.Client {
	addr := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword, // "" nếu không có password
		DB:       cfg.RedisDB,       // 0 mặc định
	})

	// test kết nối
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to Redis at %s: %v", addr, err)
	}

	return rdb
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
