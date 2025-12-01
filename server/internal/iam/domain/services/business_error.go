package services

import "server/pkg/errors"

type BusinessError interface {
	// WithMessage(msg string) func(*errors.BusinessError)
	// WithData(data map[string]interface{}) func(*errors.BusinessError)
	GetErrorByCode(code string) *errors.BusinessError
	// NewBusinessError(template errors.BusinessError, opts ...func(*errors.BusinessError)) *errors.BusinessError
	// GetBusinessError(err error) *errors.BusinessError
}
