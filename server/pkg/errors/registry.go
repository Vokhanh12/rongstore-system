package errors

func New(template AppError, opts ...func(*AppError)) *AppError {
	e := copy(template)
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func WithMessage(msg string) func(*AppError) {
	return func(e *AppError) { e.Message = msg }
}

func WithData(data map[string]interface{}) func(*AppError) {
	return func(e *AppError) { e.Data = data }
}

func copy(src AppError) *AppError {
	dst := src
	if src.Data != nil {
		dst.Data = make(map[string]interface{}, len(src.Data))
		for k, v := range src.Data {
			dst.Data[k] = v
		}
	}
	return &dst
}
