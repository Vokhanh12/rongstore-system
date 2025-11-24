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

type Keycloak interface {
	GetToken(ctx context.Context, username, password string) (*Token, *errors.BusinessError)
	RefreshToken(ctx context.Context, refreshToken string) (*Token, *errors.BusinessError)
	CheckHealth() *errors.BusinessError
	GetBaseURL() string
	//GetUserPermissions(ctx context.Context, accessToken string) ([]Permission, *errors.BusinessError)
}
