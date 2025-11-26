package error

import (
	sv "server/internal/iam/domain/services"
	"server/pkg/errors"
)

var _ sv.BusinessError = (*IamBusinessError)(nil)

type IamBusinessError struct {
	defaultArrErrs map[string]errors.BusinessError
}

// Clone implements services.BusinessError.
func (i *IamBusinessError) Clone(be sv.BusinessError) *sv.BusinessError {
	c := be
	return &c
}

// GetBusinessError implements services.BusinessError.
func (i *IamBusinessError) GetBusinessError(arrErrs map[string]sv.BusinessError, err error) *sv.BusinessError {
	if be, ok := err.(*BusinessError); ok {
		return be
	}
	return i.Clone(UNKNOWN_DOMAIN_KEY)
}

// NewBusinessError implements services.BusinessError.
func (i *IamBusinessError) NewBusinessError(template sv.BusinessError, opts ...func(*sv.BusinessError)) *sv.BusinessError {
	be := template
	for _, opt := range opts {
		opt(&be)
	}
	return &be
}

// GetErrorByCode implements services.BusinessError.
func (i *IamBusinessError) GetErrorByCode(code string) *sv.BusinessError {
	if err, ok := i.defaultArrErrs[code]; ok {
		return &err
	}
	return i.Clone(UNKNOWN_DOMAIN_KEY)
}

func InitIamBusinessError(arrErrs map[string]errors.BusinessError) sv.BusinessError {
	return &IamBusinessError{
		defaultArrErrs: arrErrs,
	}
}

func (i *IamBusinessError) WithData(data map[string]interface{}) {
	panic("unimplemented")
}

func (i *IamBusinessError) WithMessage(msg string) func(*sv.BusinessError) {
	panic("unimplemented")
}
