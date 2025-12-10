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

	domain_errors "server/internal/iam/domain"
	sv "server/internal/iam/domain/services"
	"server/pkg/config"
	"server/pkg/errors"
	"server/pkg/logger"
	"server/pkg/trace"
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

func InitKeycloakClient(ctx context.Context, cfg *config.Config) sv.Keycloak {
	maxRetries := cfg.MaxRetries
	interval := time.Duration(cfg.Interval) * time.Second

	kc := NewKeycloakClient(cfg, infbe)

	for i := 0; i < maxRetries; i++ {
		if err := kc.CheckHealth(ctx); err != nil {
			return kc
		} else {

			fields := map[string]interface{}{
				"trace_id":  trace.NewTraceID(),
				"retry":     i + 1,
				"max_retry": maxRetries,
				"operation": "init.keycloak.client",
			}

			if i < maxRetries-1 {
				logger.LogInfraDebug(ctx, err, "", fields)
			} else {
				logger.LogBySeverity(ctx, err, fields)
			}
		}

		time.Sleep(interval * time.Duration(1<<i))
	}

	be := domain_errors.KEYCLOAK_UNAVAILABLE
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

func (kc *KeycloakClient) CheckHealth(ctx context.Context) error {
	resp, err := kc.Client.Get(kc.Health)

	if err != nil {
		fields := map[string]interface{}{
			"trace_id":  trace.NewTraceID(),
			"operation": "keycloak.checkhealth",
		}
		logger.LogBySeverity(ctx, err, fields)
		return errors.Clone(domain_errors.KEYCLOAK_UNAVAILABLE)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {

	case http.StatusOK:
		return nil

	case http.StatusUnauthorized, http.StatusForbidden:
		return errors.Clone(domain_errors.KEYCLOAK_CONFIG_INVALID)

	case http.StatusNotFound:
		return errors.Clone(domain_errors.KEYCLOAK_HEALTH_ENDPOINT_INVALID)

	case http.StatusInternalServerError:
		return errors.Clone(domain_errors.KEYCLOAK_INTERNAL)

	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return errors.Clone(domain_errors.KEYCLOAK_UNAVAILABLE)

	default:
		return errors.Clone(domain_errors.KEYCLOAK_UNAVAILABLE)
	}
}

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
		return nil, errors.Clone(domain_errors.KEYCLOAK_UNAVAILABLE)
	}

	if status != http.StatusOK {
		return kc.mapKeycloakError(respBody)
	}

	var token sv.Token
	if err := json.Unmarshal(respBody, &token); err != nil {
		return nil, errors.Clone(errors.INTERNAL_FALLBACK)
	}

	if token.AccessToken == "" {
		return nil, errors.Clone(domain_errors.KEYCLOAK_UNAVAILABLE)
	}

	return &token, nil
}

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
		return nil, errors.Clone(domain_errors.KEYCLOAK_UNAVAILABLE)
	}

	if status != http.StatusOK {
		return kc.mapKeycloakError(respBody)
	}

	var token sv.Token
	if err := json.Unmarshal(respBody, &token); err != nil {
		return nil, errors.Clone(errors.INTERNAL_FALLBACK)
	}

	if token.AccessToken == "" {
		return nil, errors.Clone(domain_errors.KEYCLOAK_UNAVAILABLE)
	}

	return &token, nil
}

func (kc *KeycloakClient) mapKeycloakError(respBody []byte) (*sv.Token, *errors.BusinessError) {
	var kcErr struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}

	if err := json.Unmarshal(respBody, &kcErr); err != nil {
		return nil, errors.Clone(errors.INTERNAL_FALLBACK)
	}

	switch kcErr.Error {
	case "invalid_grant":
		return nil, errors.Clone(domain_errors.INVALID_CREDENTIALS)
	default:
		return nil, errors.Clone(domain_errors.KEYCLOAK_UNAVAILABLE)
	}
}
