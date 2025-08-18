package errors

import (
	"context"
	"errors"
	"net/http"

	dom "server/internal/iam/domain"
	"server/pkg/metrics"

	"go.uber.org/zap"
)

// TranslateDomainError converts any error into an *AppError, instruments metrics and logs.
// Note: 'service' and 'handler' labels are set to "iam" and "unknown" here. If you want
// more precise handler labels, consider storing handler info in context.
func TranslateDomainError(ctx context.Context, err error) *AppError {
	if err == nil {
		return nil
	}

	// if already AppError, normalize and return
	var appErr *AppError
	if errors.As(err, &appErr) {
		if appErr.TraceID == "" {
			appErr.TraceID = TraceIDFromContext(ctx)
		}
		if appErr.Status == 0 {
			appErr.Status = http.StatusInternalServerError
		}
		// instrument & log
		incAndLog(ctx, appErr, err)
		return appErr
	}

	// domain business error?
	var be *dom.BusinessError
	if errors.As(err, &be) {
		var mapped mapping
		var ok bool

		// 1) try YAML loader mapping
		if mapped, ok = lookupDomain(be.Key); ok {
			appErr = New(mapped.Code, mapped.Status, mapped.Message, TraceIDFromContext(ctx), err)
			incAndLog(ctx, appErr, err)
			return appErr
		}

		// 2) unknown domain key -> default unknown (from YAML defaults or hardcoded)
		unk := defaultUnknown()
		msg := be.Message
		if msg == "" {
			msg = unk.Message
		}
		appErr = New(unk.Code, unk.Status, msg, TraceIDFromContext(ctx), err)
		incAndLog(ctx, appErr, err)
		return appErr
	}

	// 3) non-domain error -> internal fallback
	df := defaultInternal()
	appErr = New(df.Code, df.Status, df.Message, TraceIDFromContext(ctx), err)
	incAndLog(ctx, appErr, err)
	return appErr
}

// incAndLog increments the error metric and logs the translated error.
// Uses service="iam" and handler="unknown" by default.
func incAndLog(ctx context.Context, appErr *AppError, originalErr error) {
	// safe-guard: metrics package must be registered at startup
	if metrics.ErrorsTotal != nil {
		// labels: service, handler, code
		metrics.ErrorsTotal.WithLabelValues("iam", "unknown", appErr.Code).Inc()
	}

	// log: use Warn for client/domain errors (<500), Error for server errors (>=500)
	trace := TraceIDFromContext(ctx)
	if appErr.Status >= 500 {
		zap.L().Error("translated_error",
			zap.String("service", "iam"),
			zap.String("handler", "unknown"),
			zap.String("code", appErr.Code),
			zap.Int("status", appErr.Status),
			zap.String("trace_id", trace),
			zap.Error(originalErr),
		)
	} else {
		zap.L().Warn("translated_error",
			zap.String("service", "iam"),
			zap.String("handler", "unknown"),
			zap.String("code", appErr.Code),
			zap.Int("status", appErr.Status),
			zap.String("trace_id", trace),
			zap.Error(originalErr),
		)
	}
}
