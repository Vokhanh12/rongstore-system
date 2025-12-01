package errors

import (
	sterrors "errors"
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

func Clone(be BusinessError) *BusinessError {
	newBe := be
	if be.Data != nil {
		newData := make(map[string]interface{}, len(be.Data))
		for k, v := range be.Data {
			newData[k] = v
		}
		newBe.Data = newData
	}
	return &newBe
}

func GetBusinessError(err error) *BusinessError {
	if err == nil {
		return nil
	}

	var be *BusinessError
	if sterrors.As(err, &be) {
		return be
	}

	type coder interface {
		Code() string
	}

	var c coder
	if sterrors.As(err, &c) {
		if mapped, ok := ErrorByCode[c.Code()]; ok {
			be := mapped
			return &be
		}
	}

	return Clone(UNKNOWN_DOMAIN_KEY)
}

func WithData(data map[string]interface{}) func(*BusinessError) {
	return func(be *BusinessError) {
		be.Data = data
	}
}

func WithMessage(msg string) func(*BusinessError) {
	return func(be *BusinessError) {
		be.Message = msg
	}
}
