package services

import (
	"context"
	"server/pkg/errors"
)

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

type Permission struct {
	Rsid   string   `json:"rsid"`
	Rsname string   `json:"rsname"`
	Scopes []string `json:"scopes"`
}

type IntrospectionResult struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
	Exp       int64  `json:"exp"`
	Iat       int64  `json:"iat"`
	Nbf       int64  `json:"nbf"`
	Sub       string `json:"sub"`
	Aud       string `json:"aud"`
	Iss       string `json:"iss"`
	Jti       string `json:"jti"`
}

type Keycloak interface {
	// AUTHENTICATION
	GetToken(ctx context.Context, username, password string) (*Token, *errors.AppError)
	RefreshToken(ctx context.Context, refreshToken string) (*Token, *errors.AppError)
	Logout(ctx context.Context, refreshToken string) *errors.AppError

	// TOKEN UTILITIES
	IntrospectToken(ctx context.Context, token string) (*IntrospectionResult, *errors.AppError)

	// AUTHORIZATION (OPTIONAL)
	GetUserPermissions(ctx context.Context, accessToken string) ([]Permission, *errors.AppError)

	// SERVICE HEALTH
	CheckHealth(ctx context.Context) *errors.AppError
	GetBaseURL() string
}
