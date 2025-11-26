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
	"server/pkg/logger"
)

var _ sv.Keycloak = (*KeycloakClient)(nil)

type KeycloakClient struct {
	BaseURL string
	Client  *http.Client
	Config  *config.Config
	Health  string
}

// func (kc *KeycloakClient) GetUserPermissions(ctx context.Context, accessToken string) ([]sv.Permission, *errors.BusinessError) {

// 	form := url.Values{}
// 	form.Set("grant_type", kc.Config.KeycloakGrantUmaTicketType)
// 	form.Set("audience", kc.Config.KeycloakAudience)
// 	form.Set("response_mode", kc.Config.KeycloakResponsePermissionsMode)

// 	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", kc.BaseURL, kc.Config.KeycloakRealm)
// 	body := bytes.NewBufferString(form.Encode())

// 	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)

// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	return nil, nil

// }
// =====================================================
// INIT — giống hệt RedisSessionStore pattern
// =====================================================
func InitKeycloakClient(ctx context.Context, cfg *config.Config) sv.Keycloak {
	maxRetries := cfg.MaxRetries
	interval := time.Duration(cfg.Interval) * time.Second

	kc := NewKeycloakClient(cfg)

	for i := 0; i < maxRetries; i++ {
		if err := kc.CheckHealth(); err == nil {
			return kc
		} else {
			be := businessError.GetBusinessError(err)
			fields := map[string]interface{}{
				"retry":     i + 1,
				"operation": "init.keycloak.client",
				"error":     err.Error(),
			}

			logger.LogBySeverity(ctx, *be, fields)
		}

		time.Sleep(interval)
	}

	be := businessError.KEYCLOAK_UNAVAILABLE
	panic(fmt.Sprintf(
		"PANIC: [%s][%s] %s | cause: %s | server_action: %s | retryable: %v",
		be.Code,
		be.Key,
		be.Message,
		be.Cause,
		be.ServerAction,
		be.Retryable,
	))
}

// =====================================================
// CONSTRUCTOR
// =====================================================
func NewKeycloakClient(cfg *config.Config) sv.Keycloak {
	return &KeycloakClient{
		BaseURL: cfg.KeycloakURL,
		Client:  &http.Client{Timeout: 5 * time.Second},
		Config:  cfg,
		Health:  cfg.KeycloakServerHealth,
	}
}

func (kc *KeycloakClient) GetBaseURL() string {
	return kc.BaseURL
}

// =====================================================
// HEALTH CHECK
// =====================================================
func (kc *KeycloakClient) CheckHealth() *errors.BusinessError {
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

// =====================================================
// TOKEN REQUEST HELPERS
// =====================================================
func (kc *KeycloakClient) tokenEndpoint() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		kc.BaseURL,
		kc.Config.KeycloakRealm,
	)
}

func (kc *KeycloakClient) doFormRequest(
	ctx context.Context,
	form url.Values,
) ([]byte, int, error) {

	url := kc.tokenEndpoint()
	body := bytes.NewBufferString(form.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := kc.Client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	respBody, readErr := io.ReadAll(resp.Body)
	return respBody, resp.StatusCode, readErr
}

// =====================================================
// GET TOKEN
// =====================================================
func (kc *KeycloakClient) GetToken(
	ctx context.Context,
	username, password string,
) (*sv.Token, *errors.BusinessError) {

	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("client_id", kc.Config.KeycloakClientID)
	form.Set("client_secret", kc.Config.KeycloakSecret)
	form.Set("username", username)
	form.Set("password", password)
	form.Set("scope", kc.Config.KeycloakScope)

	respBody, status, err := kc.doFormRequest(ctx, form)
	if err != nil {
		return nil, errors.Clone(businessError.KEYCLOAK_UNAVAILABLE)
	}

	if status != http.StatusOK {
		return kc.mapKeycloakError(respBody)
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

// =====================================================
// REFRESH TOKEN
// =====================================================
func (kc *KeycloakClient) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*sv.Token, *errors.BusinessError) {

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", kc.Config.KeycloakClientID)
	form.Set("client_secret", kc.Config.KeycloakSecret)
	form.Set("refresh_token", refreshToken)

	respBody, status, err := kc.doFormRequest(ctx, form)
	if err != nil {
		return nil, errors.Clone(businessError.KEYCLOAK_UNAVAILABLE)
	}

	if status != http.StatusOK {
		return kc.mapKeycloakError(respBody)
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

// =====================================================
// ERROR MAPPING
// =====================================================
func (kc *KeycloakClient) mapKeycloakError(respBody []byte) (*sv.Token, *errors.BusinessError) {
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
