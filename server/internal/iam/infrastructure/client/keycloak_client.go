package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	sv "server/internal/iam/domain/services"
	"server/pkg/config"
)

var _ sv.Keycloak = (*KeycloakClient)(nil)

type KeycloakClient struct {
	BaseURL string
	Client  *http.Client
	Config  *config.Config
	Health  string
}

func asStdContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	if std, ok := ctx.(context.Context); ok {
		return std
	}
	return context.Background()
}

func (kc *KeycloakClient) GetToken(ctx context.Context, username, password string) (*sv.Token, error) {
	stdCtx := asStdContext(ctx)

	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("client_id", kc.Config.KeycloakClientID)
	form.Set("client_secret", kc.Config.KeycloakSecret)
	form.Set("username", username)
	form.Set("password", password)
	form.Set("scope", kc.Config.KeycloakScope)

	req, err := http.NewRequestWithContext(stdCtx, http.MethodPost,
		fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", kc.BaseURL, kc.Config.KeycloakRealm),
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := kc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get token: status %d, body %s", resp.StatusCode, string(body))
	}

	var token sv.Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (kc *KeycloakClient) RefreshToken(ctx context.Context, refreshToken string) (*sv.Token, error) {
	stdCtx := asStdContext(ctx)

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", kc.Config.KeycloakClientID)
	form.Set("client_secret", kc.Config.KeycloakSecret)
	form.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(stdCtx, http.MethodPost,
		fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", kc.BaseURL, kc.Config.KeycloakRealm),
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := kc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to refresh token: status %d, body %s", resp.StatusCode, string(body))
	}

	var token sv.Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (kc *KeycloakClient) CheckHealth() error {
	resp, err := kc.Client.Get(kc.Health)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Keycloak health check failed: status %d, body %s", resp.StatusCode, string(body))
	}
	return nil
}

func (kc *KeycloakClient) GetBaseURL() string {
	return kc.BaseURL
}

func NewKeycloakClient(cfg *config.Config) sv.Keycloak {
	return &KeycloakClient{
		BaseURL: cfg.KeycloakURL,
		Client:  &http.Client{Timeout: 5 * time.Second},
		Config:  cfg,
		Health:  cfg.KeycloakServerHealth,
	}
}

func InitKeycloakClient(cfg *config.Config, maxRetries int, interval time.Duration) (sv.Keycloak, error) {
	kc := NewKeycloakClient(cfg)
	for i := 0; i < maxRetries; i++ {
		if err := kc.CheckHealth(); err == nil {
			return kc, nil
		}
		time.Sleep(interval)
	}
	return nil, fmt.Errorf("Keycloak not ready after %d retries", maxRetries)
}
