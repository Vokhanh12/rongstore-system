package services

import (
	"context"
	"time"
)

type SessionEntry struct {
	SessionID string
	ClientPub []byte
	ServerPub []byte

	Kc2s []byte
	Ks2c []byte

	HKDFSalt []byte
	Expiry   time.Time

	UserID   string
	ClientIP string
}

// SessionCache
type IRedisCache interface {
	//	CheckHealth(r *RedisSessionStore) *errors.BusinessError
	StoreSession(ctx context.Context, e *SessionEntry) error
	GetSession(ctx context.Context, sessionID string) (*SessionEntry, error)
	DeleteSession(ctx context.Context, sessionID string) error
	CheckAndRecordNonceAtomic(ctx context.Context, sessionID, nonceB64 string, windowSeconds int) (bool, error)
}

// RedisNonceService
type INonceCache interface {
	CheckAndRecordNonceAtomic(ctx context.Context, sessionID, nonceB64 string, windowSeconds int) (bool, error)
}

// RedisHealth
type ICacheHealth interface {
	Ping(ctx context.Context) error
}

// Rate limit
type IRateLimiter interface {
	IncreaseLoginAttempt(ctx context.Context, userID string) (int, error)
	ResetLoginAttempt(ctx context.Context, userID string) error
}
