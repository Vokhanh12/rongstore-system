package cache

import (
	"context"
	"encoding/base64"
	"strconv"
	"time"

	"server/internal/iam/domain/services"
	"server/pkg/config"

	"github.com/redis/go-redis/v9"
)

type RedisSessionStore struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRedisSessionStore(rdb *redis.Client, ttl time.Duration) *RedisSessionStore {
	return &RedisSessionStore{rdb: rdb, ttl: ttl}
}

func RedisTTLFromConfig(cfg *config.Config) time.Duration {
	return time.Duration(cfg.RedisTTL) * time.Second
}

func sessionKey(sessionID string) string {
	return "session:" + sessionID
}

func (r *RedisSessionStore) StoreSession(ctx context.Context, e *services.SessionEntry) error {
	key := sessionKey(e.SessionID)
	fields := map[string]interface{}{
		"client_pub": base64.StdEncoding.EncodeToString(e.ClientPub),
		"server_pub": base64.StdEncoding.EncodeToString(e.ServerPub),
		"kc2s":       base64.StdEncoding.EncodeToString(e.Kc2s),
		"ks2c":       base64.StdEncoding.EncodeToString(e.Ks2c),
		"hkdf_salt":  base64.StdEncoding.EncodeToString(e.HKDFSalt),
		"expiry":     e.Expiry.Unix(),
	}

	if err := r.rdb.HSet(ctx, key, fields).Err(); err != nil {
		return err
	}
	return r.rdb.Expire(ctx, key, r.ttl).Err()
}

func (r *RedisSessionStore) GetSession(ctx context.Context, sessionID string) (*services.SessionEntry, error) {
	key := sessionKey(sessionID)
	m, err := r.rdb.HGetAll(ctx, key).Result()
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

	return &services.SessionEntry{
		SessionID: sessionID,
		ClientPub: clientPub,
		ServerPub: serverPub,
		Kc2s:      kc2s,
		Ks2c:      ks2c,
		HKDFSalt:  hkdfSalt,
		Expiry:    exp,
	}, nil
}

func (r *RedisSessionStore) DeleteSession(ctx context.Context, sessionID string) error {
	key := sessionKey(sessionID)
	return r.rdb.Del(ctx, key).Err()
}

func (r *RedisSessionStore) CheckAndRecordNonceAtomic(ctx context.Context, sessionID, nonceB64 string, windowSeconds int) (bool, error) {
	nonceKey := sessionKey(sessionID) + ":nonce:" + nonceB64
	return r.rdb.SetNX(ctx, nonceKey, "1", time.Duration(windowSeconds)*time.Second).Result()
}
