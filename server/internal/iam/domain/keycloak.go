package domain

import "context"

type Token struct {
	IdToken          string `json:"token_id"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int32  `json:"expires_in"`
	ReFreshExpiresIn int32  `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int32  `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type Keycloak interface {
	GetToken(ctx context.Context, username, password string) (*Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (*Token, error)
	CheckHealth() error
	GetBaseURL() string
}
