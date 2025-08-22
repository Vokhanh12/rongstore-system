package observability

import "context"

// SessionInfo is intentionally minimal and contains no secrets.
type SessionInfo struct {
	UserID string
	IP     string
	// add non-sensitive metadata if needed
}

type SessionStore interface {
	// Get returns (nil, nil) when session not found.
	Get(ctx context.Context, sessionID string) (*SessionInfo, error)
}
