package domain

import "time"

// SessionEntry lưu metadata và session keys (sensitive) cho 1 session.
// Lưu ý: Kc2s/Ks2c nên được lưu dưới dạng []byte khi dùng trong memory,
// hoặc base64 khi serialize vào Redis. Ở domain ta giữ []byte để dễ dùng.
type SessionEntry struct {
	SessionID string
	ClientPub []byte // raw bytes of client public key
	ServerPub []byte // raw bytes of server public key

	Kc2s []byte // client->server AEAD key (32 bytes)
	Ks2c []byte // server->client AEAD key (32 bytes)

	HKDFSalt []byte // optional: salt used for HKDF (not secret but handy)
	Expiry   time.Time

	// optional fields for replay protection etc can be kept in infrastructure layer
}

// SessionStore là interface domain cho session persistence.
// Implementer (Redis, in-memory, etc) phải implement các method này.
type SessionStore interface {
	// StoreSession persists a session entry and sets TTL (expiry).
	StoreSession(ctx Context, e *SessionEntry) error

	// GetSession returns session entry or error; return nil,nil if not found.
	GetSession(ctx Context, sessionID string) (*SessionEntry, error)

	// DeleteSession removes session entry and ensures sensitive bytes are zeroed if needed.
	DeleteSession(ctx Context, sessionID string) error

	// CheckAndRecordNonceAtomic checks whether nonce seen; if not, records it with TTL.
	// Returns true if nonce was not seen (accepted), false if replay (already seen).
	CheckAndRecordNonceAtomic(ctx Context, sessionID, nonceB64 string, windowSeconds int) (bool, error)
}

// Note: Domain should not import redis. Context alias to avoid importing "context" in domain file
type Context = interface{}
