package domain

// Token model
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

// Keycloak interface (domain-level)
type Keycloak interface {
	GetToken(ctx Context, username, password string) (*Token, error)
	RefreshToken(ctx Context, refreshToken string) (*Token, error)
	CheckHealth() error
	GetBaseURL() string
}
