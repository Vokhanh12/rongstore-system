package infahelper

import (
	"server/pkg/errors"
	"time"
)

func Retry[T any](
	maxRetries int,
	interval time.Duration,
	fn func() (T, *errors.AppError),
) (T, *errors.AppError) {
	var zero T
	var lastErr *errors.AppError

	for i := 0; i < maxRetries; i++ {
		res, err := fn()
		if err == nil {
			return res, nil
		}

		lastErr = err

		if !err.Retryable {
			break
		}

		time.Sleep(interval * time.Duration(1<<i))
	}

	return zero, lastErr
}
