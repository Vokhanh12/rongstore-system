package errors

import (
	"server/internal/iam/domain/services"
	"server/pkg/config"
	"server/pkg/errors"
)

var _ services.BusinessError = (*BusinessError)(nil)

type BusinessError struct {
	defaultArrErrs map[string]errors.BusinessError
}

func InitBusinessError(cfg *config.Config) services.BusinessError {
	return &BusinessError{
		defaultArrErrs: cfg.IAMDefaultErrors,
	}
}

// func (i *BusinessError) NewBusinessError(template errors.BusinessError, opts ...func(*errors.BusinessError)) *errors.BusinessError {
// 	be := template
// 	for _, opt := range opts {
// 		opt(&be)
// 	}
// 	return &be
// }

// func (i *BusinessError) GetBusinessError(err error) *errors.BusinessError {
// 	if err == nil {
// 		return nil
// 	}

// 	var be *errors.BusinessError
// 	if sterrors.As(err, &be) {
// 		return be
// 	}

// 	type coder interface {
// 		Code() string
// 	}

// 	var c coder
// 	if sterrors.As(err, &c) {
// 		if mapped, ok := errors.ErrorByCode[c.Code()]; ok {
// 			be := mapped
// 			return &be
// 		}
// 	}

// 	return errors.Clone(errors.UNKNOWN_DOMAIN_KEY)
// }

func (i *BusinessError) GetErrorByCode(code string) *errors.BusinessError {
	if err, ok := i.defaultArrErrs[code]; ok {
		return &err
	}
	return errors.Clone(errors.UNKNOWN_DOMAIN_KEY)
}

// func (i *BusinessError) WithData(data map[string]interface{}) func(*errors.BusinessError) {
// 	return func(be *errors.BusinessError) {
// 		be.Data = data
// 	}
// }

// func (i *BusinessError) WithMessage(msg string) func(*errors.BusinessError) {
// 	return func(be *errors.BusinessError) {
// 		be.Message = msg
// 	}
// }
