package error

import (
	"server/internal/iam/domain/services"
	"server/pkg/errors"
)

var _ services.BusinessError = (*IamBusinessError)(nil)

type IamBusinessError struct {
	defaultArrErrs map[string]errors.BusinessError
}

func InitIamBusinessError(arrErrs map[string]errors.BusinessError) *IamBusinessError {
	return &IamBusinessError{
		defaultArrErrs: arrErrs,
	}
}

func (i *IamBusinessError) NewBusinessError(template errors.BusinessError, opts ...func(*errors.BusinessError)) *errors.BusinessError {
	be := template
	for _, opt := range opts {
		opt(&be)
	}
	return &be
}

func (i *IamBusinessError) GetBusinessError(arrErrs map[string]errors.BusinessError, err error) *errors.BusinessError {
	if be, ok := err.(*errors.BusinessError); ok {
		return be
	}
	return errors.Clone(errors.UNKNOWN_DOMAIN_KEY)
}

func (i *IamBusinessError) GetErrorByCode(code string) *errors.BusinessError {
	if err, ok := i.defaultArrErrs[code]; ok {
		return &err
	}
	return errors.Clone(errors.UNKNOWN_DOMAIN_KEY)
}

func (i *IamBusinessError) WithData(data map[string]interface{}) func(*errors.BusinessError) {
	return func(be *errors.BusinessError) {
		be.Data = data
	}
}

func (i *IamBusinessError) WithMessage(msg string) func(*errors.BusinessError) {
	return func(be *errors.BusinessError) {
		be.Message = msg
	}
}
