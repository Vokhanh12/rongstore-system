package domain

import "time"

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

type SessionStore interface {
	StoreSession(ctx Context, e *SessionEntry) error
	GetSession(ctx Context, sessionID string) (*SessionEntry, error)
	DeleteSession(ctx Context, sessionID string) error
	CheckAndRecordNonceAtomic(ctx Context, sessionID, nonceB64 string, windowSeconds int) (bool, error)
}
