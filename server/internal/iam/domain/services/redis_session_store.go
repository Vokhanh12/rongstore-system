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

type RedisSessionStore interface {
	//	CheckHealth(r *RedisSessionStore) *errors.BusinessError
	StoreSession(ctx context.Context, e *SessionEntry) error
	GetSession(ctx context.Context, sessionID string) (*SessionEntry, error)
	DeleteSession(ctx context.Context, sessionID string) error
	CheckAndRecordNonceAtomic(ctx context.Context, sessionID, nonceB64 string, windowSeconds int) (bool, error)
}
