// cache/obs_adapter.go
package cache

import (
	"context"
	"server/pkg/observability"
)

// redisObsAdapter wraps RedisSessionStore to observability.SessionStore
type redisObsAdapter struct {
	s *RedisSessionStore
}

func NewObservabilitySessionStoreAdapter(s *RedisSessionStore) observability.SessionStore {
	return &redisObsAdapter{s: s}
}

func (a *redisObsAdapter) Get(ctx context.Context, sessionID string) (*observability.SessionInfo, error) {
	se, err := a.s.GetSession(ctx, sessionID)
	if err != nil || se == nil {
		return nil, err
	}

	return &observability.SessionInfo{
		UserID: se.UserID,
		IP:     se.ClientIP,
	}, nil
}
