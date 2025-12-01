package cache

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"server/internal/iam/domain/services"
	sv "server/internal/iam/domain/services"
	"server/pkg/config"
	"server/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type RedisSessionStore struct {
	rdb *redis.Client
	ttl time.Duration
}

func InitRedisSessionStore(ctx context.Context, cfg *config.Config, infbe services.BusinessError) sv.RedisSessionStore {
	maxRetries := cfg.MaxRetries
	interval := time.Duration(cfg.Interval) * time.Second

	for i := 0; i < maxRetries; i++ {
		rdb := newRedisClient(cfg)

		err := pingRedis(rdb)
		if err == nil {
			return &RedisSessionStore{
				rdb: rdb,
				ttl: RedisTTLFromConfig(cfg),
			}
		}

		fields := map[string]interface{}{
			"retry":     i + 1,
			"operation": "init.redis.session.store",
		}

		if i < maxRetries-1 {
			logger.LogInfraDebug(ctx, err, "", fields)
		} else {
			logger.LogBySeverity(ctx, err, fields)
		}

		time.Sleep(interval * time.Duration(1<<i))
	}

	return nil
}

func newRedisClient(cfg *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
}

func pingRedis(rdb *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return rdb.Ping(ctx).Err()
}

// func (r *RedisSessionStore) CheckHealth() *errors.BusinessError {
// 	if err := pingRedis(r.rdb); err != nil {
// 		return &domain_errors.REDIS_UNAVAILABLE
// 	}
// 	return nil
// }

func RedisTTLFromConfig(cfg *config.Config) time.Duration {
	return time.Duration(cfg.RedisTTL) * time.Second
}

func sessionKey(sessionID string) string {
	return "session:" + sessionID
}

func (r *RedisSessionStore) StoreSession(ctx context.Context, e *sv.SessionEntry) error {
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

func (r *RedisSessionStore) GetSession(ctx context.Context, sessionID string) (*sv.SessionEntry, error) {
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

	return &sv.SessionEntry{
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
