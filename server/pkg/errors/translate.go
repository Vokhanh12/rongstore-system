package errors

import (
	"context"
	"errors"

	dom "server/internal/iam/domain"
)

func TranslateDomainError(ctx context.Context, err error) *AppError {
	if err == nil {
		return nil
	}

	var be *dom.BusinessError
	if errors.As(err, &be) {
		if m, ok := handshakeMap[be.Key]; ok {
			return New(m.Code, m.Status, m.Message, TraceIDFromContext(ctx), err)
		}
		// Có thể thêm map cho các usecase khác ở đây
	}

	// fallback generic
	return New("INTERNAL-500", 500, "Internal server error", TraceIDFromContext(ctx), err)
}
