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
	"server/pkg/util/infahelper"
)

var _ sv.Keycloak = (*KeycloakClient)(nil)

type KeycloakClient struct {
	BaseURL string
	Client  *http.Client
	Config  *config.Config
	Health  string
}

// GetUserPermissions implements services.Keycloak.
func (kc *KeycloakClient) GetUserPermissions(ctx context.Context, accessToken string) ([]sv.Permission, *errors.AppError) {
	panic("unimplemented")
}

// IntrospectToken implements services.Keycloak.
func (kc *KeycloakClient) IntrospectToken(ctx context.Context, token string) (*sv.IntrospectionResult, *errors.AppError) {
	panic("unimplemented")
}

// Logout implements services.Keycloak.
func (kc *KeycloakClient) Logout(ctx context.Context, refreshToken string) *errors.AppError {
	panic("unimplemented")
}

// func (kc *KeycloakClient) GetUserPermissions(ctx context.Context, accessToken string) ([]sv.Permission, *errors.AppError) {

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

func InitKeycloakClient(
	ctx context.Context,
	cfg *config.Config,
) sv.Keycloak {

	kc, err := infahelper.Retry(
		cfg.MaxRetries,
		time.Duration(cfg.Interval)*time.Second,
		func() (*KeycloakClient, *errors.AppError) {
			client := NewKeycloakClient(cfg)
			if err := client.CheckHealth(ctx); err != nil {
				return nil, err
			}
			return client, nil
		},
	)

	if err != nil {
		logger.LogBySeverity(
			ctx,
			"Init.KeycloakClient",
			err,
		)
		return nil
	}

	return kc
}

func NewKeycloakClient(cfg *config.Config) *KeycloakClient {
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

func (kc *KeycloakClient) CheckHealth(ctx context.Context) *errors.AppError {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		kc.Health,
		nil,
	)
	if err != nil {
		return errors.New(
			domain_errors.KEYCLOAK_UNAVAILABLE,
			errors.SetError(err),
		)
	}

	resp, err := kc.Client.Do(req)
	if err != nil {
		return errors.New(
			domain_errors.KEYCLOAK_UNAVAILABLE,
			errors.SetError(err),
		)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil

	case http.StatusUnauthorized, http.StatusForbidden:
		return errors.New(domain_errors.KEYCLOAK_CONFIG_INVALID)

	case http.StatusNotFound:
		return errors.New(domain_errors.KEYCLOAK_HEALTH_ENDPOINT_INVALID)

	case http.StatusInternalServerError:
		return errors.New(domain_errors.KEYCLOAK_INTERNAL)

	case http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return errors.New(domain_errors.KEYCLOAK_UNAVAILABLE)

	default:
		return errors.New(domain_errors.KEYCLOAK_UNAVAILABLE)
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
) (*sv.Token, *errors.AppError) {

	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("client_id", kc.Config.KeycloakClientID)
	form.Set("client_secret", kc.Config.KeycloakSecret)
	form.Set("username", username)
	form.Set("password", password)
	form.Set("scope", kc.Config.KeycloakScope)

	respBody, status, err := kc.doFormRequest(ctx, form)
	if err != nil {
		return nil, errors.New(
			domain_errors.KEYCLOAK_UNAVAILABLE,
			errors.SetError(err))
	}

	if status != http.StatusOK {
		return kc.mapKeycloakError(respBody)
	}

	var token sv.Token
	if err := json.Unmarshal(respBody, &token); err != nil {
		return nil, errors.New(
			errors.INTERNAL_FALLBACK,
			errors.SetError(err))
	}

	if token.AccessToken == "" {
		return nil, errors.New(domain_errors.KEYCLOAK_UNAVAILABLE)
	}

	return &token, nil
}

func (kc *KeycloakClient) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*sv.Token, *errors.AppError) {

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", kc.Config.KeycloakClientID)
	form.Set("client_secret", kc.Config.KeycloakSecret)
	form.Set("refresh_token", refreshToken)

	respBody, status, err := kc.doFormRequest(ctx, form)
	if err != nil {
		return nil, errors.New(domain_errors.KEYCLOAK_UNAVAILABLE)
	}

	if status != http.StatusOK {
		return kc.mapKeycloakError(respBody)
	}

	var token sv.Token
	if err := json.Unmarshal(respBody, &token); err != nil {
		return nil, errors.New(errors.INTERNAL_FALLBACK)
	}

	if token.AccessToken == "" {
		return nil, errors.New(domain_errors.KEYCLOAK_UNAVAILABLE)
	}

	return &token, nil
}

func (kc *KeycloakClient) mapKeycloakError(respBody []byte) (*sv.Token, *errors.AppError) {
	var kcErr struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}

	if err := json.Unmarshal(respBody, &kcErr); err != nil {
		return nil, errors.New(errors.INTERNAL_FALLBACK)
	}

	switch kcErr.Error {
	case "invalid_grant":
		return nil, errors.New(domain_errors.INVALID_CREDENTIALS)
	default:
		return nil, errors.New(domain_errors.KEYCLOAK_UNAVAILABLE)
	}
}
