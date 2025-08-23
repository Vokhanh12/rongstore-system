package cache

import (
	"context"
	"encoding/base64"
	"strconv"
	"time"

	"server/internal/iam/domain"

	"server/pkg/config"

	"github.com/redis/go-redis/v9"
)

// RedisSessionStore implements domain.SessionStore using Redis.
type RedisSessionStore struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewRedisSessionStore creates a new RedisSessionStore.
func NewRedisSessionStore(rdb *redis.Client, ttl time.Duration) domain.SessionStore {
	return &RedisSessionStore{rdb: rdb, ttl: ttl}
}

// Provider cho Wire: trả time.Duration trực tiếp
func RedisTTLFromConfig(cfg *config.Config) time.Duration {
	// giả sử cfg.RedisTTL là số giây (int)
	return time.Duration(cfg.RedisTTL) * time.Second
}

// helper convert domain.Context -> context.Context
func toStdContext(ctx domain.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	// Nếu ctx là context.Context thì return trực tiếp
	if c, ok := ctx.(context.Context); ok {
		return c
	}
	return context.Background()
}

// sessionKey helper
func sessionKey(sessionID string) string {
	return "session:" + sessionID
}

// StoreSession lưu session vào Redis hash
func (r *RedisSessionStore) StoreSession(ctx domain.Context, e *domain.SessionEntry) error {
	c := toStdContext(ctx)
	key := sessionKey(e.SessionID)
	fields := map[string]interface{}{
		"client_pub": base64.StdEncoding.EncodeToString(e.ClientPub),
		"server_pub": base64.StdEncoding.EncodeToString(e.ServerPub),
		"kc2s":       base64.StdEncoding.EncodeToString(e.Kc2s),
		"ks2c":       base64.StdEncoding.EncodeToString(e.Ks2c),
		"hkdf_salt":  base64.StdEncoding.EncodeToString(e.HKDFSalt),
		"expiry":     e.Expiry.Unix(),
	}

	if err := r.rdb.HSet(c, key, fields).Err(); err != nil {
		return err
	}
	return r.rdb.Expire(c, key, r.ttl).Err()
}

// GetSession lấy session từ Redis
func (r *RedisSessionStore) GetSession(ctx domain.Context, sessionID string) (*domain.SessionEntry, error) {
	c := toStdContext(ctx)
	key := sessionKey(sessionID)
	m, err := r.rdb.HGetAll(c, key).Result()
	if err != nil || len(m) == 0 {
		return nil, err
	}

	decode := func(k string) ([]byte, error) {
		v, ok := m[k]
		if !ok || v == "" {
			return nil, nil
		}
		return base64.StdEncoding.DecodeString(v)
	}

	clientPub, _ := decode("client_pub")
	serverPub, _ := decode("server_pub")
	kc2s, _ := decode("kc2s")
	ks2c, _ := decode("ks2c")
	hkdfSalt, _ := decode("hkdf_salt")

	exp := time.Now()
	if ts, ok := m["expiry"]; ok && ts != "" {
		if sec, err := strconv.ParseInt(ts, 10, 64); err == nil {
			exp = time.Unix(sec, 0)
		}
	}

	return &domain.SessionEntry{
		SessionID: sessionID,
		ClientPub: clientPub,
		ServerPub: serverPub,
		Kc2s:      kc2s,
		Ks2c:      ks2c,
		HKDFSalt:  hkdfSalt,
		Expiry:    exp,
	}, nil
}

// DeleteSession xóa session khỏi Redis
func (r *RedisSessionStore) DeleteSession(ctx domain.Context, sessionID string) error {
	c := toStdContext(ctx)
	key := sessionKey(sessionID)
	return r.rdb.Del(c, key).Err()
}

// CheckAndRecordNonceAtomic sử dụng SETNX để record nonce atomically
func (r *RedisSessionStore) CheckAndRecordNonceAtomic(ctx domain.Context, sessionID, nonceB64 string, windowSeconds int) (bool, error) {
	c := toStdContext(ctx)
	nonceKey := sessionKey(sessionID) + ":nonce:" + nonceB64
	return r.rdb.SetNX(c, nonceKey, "1", time.Duration(windowSeconds)*time.Second).Result()
}
