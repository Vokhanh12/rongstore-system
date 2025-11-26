package services

type BusinessError interface {
	WithMessage(msg string) func(*BusinessError)
	WithData(data map[string]interface{})
	GetErrorByCode(code string) *BusinessError
	NewBusinessError(template BusinessError, opts ...func(*BusinessError)) *BusinessError
	GetBusinessError(arrErrs map[string]BusinessError, err error) *BusinessError
	Clone(be BusinessError) *BusinessError
}
