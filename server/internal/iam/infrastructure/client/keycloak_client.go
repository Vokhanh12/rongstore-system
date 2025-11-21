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

	businessError "server/internal/iam/domain"
	sv "server/internal/iam/domain/services"
	"server/pkg/config"
	"server/pkg/errors"
)

var _ sv.Keycloak = (*KeycloakClient)(nil)

type KeycloakClient struct {
	BaseURL string
	Client  *http.Client
	Config  *config.Config
	Health  string
}

func (kc *KeycloakClient) GetUserPermissions(ctx context.Context, accessToken string) ([]sv.Permission, *errors.BusinessError) {

	form := url.Values{}
	form.Set("grant_type", kc.Config.KeycloakGrantUmaTicketType)
	form.Set("audience", kc.Config.KeycloakAudience)
	form.Set("response_mode", kc.Config.KeycloakResponsePermissionsMode)

	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", kc.BaseURL, kc.Config.KeycloakRealm)
	body := bytes.NewBufferString(form.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return nil, nil

}

func (kc *KeycloakClient) GetToken(ctx context.Context, username, password string) (*sv.Token, *errors.BusinessError) {

	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("client_id", kc.Config.KeycloakClientID)
	form.Set("client_secret", kc.Config.KeycloakSecret)
	form.Set("username", username)
	form.Set("password", password)
	form.Set("scope", kc.Config.KeycloakScope)

	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", kc.BaseURL, kc.Config.KeycloakRealm)
	body := bytes.NewBufferString(form.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)

	if err != nil {
		return nil, errors.Clone(businessError.INTERNAL_FALLBACK)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := kc.Client.Do(req)

	if err != nil {
		return nil, errors.Clone(businessError.KEYCLOAK_UNAVAILABLE)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Clone(businessError.INTERNAL_FALLBACK)
	}

	if resp.StatusCode != http.StatusOK {
		var kcErr struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if err := json.Unmarshal(respBody, &kcErr); err != nil {
			return nil, errors.Clone(businessError.INTERNAL_FALLBACK)
		}
		switch kcErr.Error {
		case "invalid_grant":
			return nil, errors.Clone(businessError.INVALID_CREDENTIALS)
		default:
			return nil, errors.Clone(businessError.KEYCLOAK_UNAVAILABLE)
		}
	}

	var token sv.Token
	if err := json.Unmarshal(respBody, &token); err != nil {
		return nil, errors.Clone(businessError.INTERNAL_FALLBACK)
	}

	if token.AccessToken == "" {
		return nil, errors.Clone(businessError.KEYCLOAK_UNAVAILABLE)
	}

	return &token, nil
}

func (kc *KeycloakClient) RefreshToken(ctx context.Context, refreshToken string) (*sv.Token, *errors.BusinessError) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", kc.Config.KeycloakClientID)
	form.Set("client_secret", kc.Config.KeycloakSecret)
	form.Set("refresh_token", refreshToken)

	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", kc.BaseURL, kc.Config.KeycloakRealm)
	body := bytes.NewBufferString(form.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Clone(businessError.INTERNAL_FALLBACK)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := kc.Client.Do(req)
	if err != nil {
		return nil, errors.Clone(businessError.KEYCLOAK_UNAVAILABLE)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Clone(businessError.INTERNAL_FALLBACK)
	}

	if resp.StatusCode != http.StatusOK {
		var kcErr struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if err := json.Unmarshal(respBody, &kcErr); err != nil {
			return nil, errors.Clone(businessError.INTERNAL_FALLBACK)
		}
		switch kcErr.Error {
		case "invalid_grant":
			return nil, errors.Clone(businessError.INVALID_CREDENTIALS)
		default:
			return nil, errors.Clone(businessError.KEYCLOAK_UNAVAILABLE)
		}
	}

	var token sv.Token
	if err := json.Unmarshal(respBody, &token); err != nil {
		return nil, errors.Clone(businessError.INTERNAL_FALLBACK)
	}

	if token.AccessToken == "" {
		return nil, errors.Clone(businessError.KEYCLOAK_UNAVAILABLE)
	}

	return &token, nil
}

func (kc *KeycloakClient) CheckHealth(ctx context.Context) (be *errors.BusinessError) {
	resp, err := kc.Client.Get(kc.Health)

	if err != nil {
		return errors.Clone(businessError.KEYCLOAK_UNAVAILABLE)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Clone(businessError.KEYCLOAK_UNAVAILABLE)
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
