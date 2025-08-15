package errors

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"trace_id,omitempty"`
	Status  int    `json:"-"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Code + ": " + e.Message
}

// New tạo AppError
func New(code string, status int, message, traceID string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		TraceID: traceID,
		Status:  status,
		Err:     err,
	}
}
