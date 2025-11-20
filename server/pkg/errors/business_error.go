package errors

import (
	"fmt"
	//iam_errors "server/internal/iam/domain"
)

type BusinessError struct {
	Code      string                 `json:"code"`
	Status    int                    `json:"status"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Severity  string                 `json:"severity,omitempty"`
	Retryable bool                   `json:"retryable,omitempty"`
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
