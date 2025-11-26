package errors

import (
	"fmt"
)

type BusinessError struct {
	Code         string                 `json:"code"`
	Status       int                    `json:"status"`
	GRPCCode     string                 `json:"grpc_code"`
	Key          string                 `json:"key"`
	Cause        string                 `json:"cause"`
	ClientAction string                 `json:"client_action"`
	ServerAction string                 `json:"server_action"`
	Message      string                 `json:"message"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Severity     string                 `json:"severity,omitempty"`
	Retryable    bool                   `json:"retryable,omitempty"`
}

func (e *BusinessError) Error() string {
	if e == nil {
		return "<nil BusinessError>"
	}
	return fmt.Sprintf("%s: %d %s", e.Code, e.Status, e.Message)
}

func NewBusinessError(template BusinessError, opts ...func(*BusinessError)) *BusinessError {
	be := template
	for _, opt := range opts {
		opt(&be)
	}
	return &be
}

func WithMessage(msg string) func(*BusinessError) {
	return func(be *BusinessError) {
		be.Message = msg
	}
}

func WithData(data map[string]interface{}) func(*BusinessError) {
	return func(be *BusinessError) {
		be.Data = data
	}
}

func (e *BusinessError) MessageOr(fallback string) string {
	if e == nil || e.Message == "" {
		return fallback
	}
	return e.Message
}

func (e *BusinessError) WithExtra(data map[string]interface{}) *BusinessError {
	if e == nil {
		return nil
	}
	if e.Data == nil {
		e.Data = map[string]interface{}{}
	}
	for k, v := range data {
		e.Data[k] = v
	}
	return e
}

func Clone(be BusinessError) *BusinessError {
	c := be
	return &c
}

func GetErrorByCode(arrErrs map[string]BusinessError, code string) *BusinessError {
	if err, ok := arrErrs[code]; ok {
		return &err
	}
	return Clone(UNKNOWN_DOMAIN_KEY)
}

func GetBusinessError(arrErrs map[string]BusinessError, err error) *BusinessError {
	if be, ok := err.(*BusinessError); ok {
		return be
	}
	return Clone(UNKNOWN_DOMAIN_KEY)
}
