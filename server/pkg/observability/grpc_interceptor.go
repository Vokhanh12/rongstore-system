package observability

import (
	"context"
	"net"
	"os"
	"strings"
	"time"

	"server/internal/iam/domain"
	"server/pkg/errors"
	"server/pkg/logger"
	"server/pkg/metrics"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// cached hostname to avoid repeated syscalls
var hostnameCache string

func init() {
	if hn, err := os.Hostname(); err == nil {
		hostnameCache = hn
	} else {
		hostnameCache = "unknown"
	}
}

// GrpcTraceUnaryInterceptor extracts "x-trace-id" / "trace-id" from metadata or generates a UUID
func GrpcTraceUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var trace string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get("x-trace-id"); len(vals) > 0 {
				trace = vals[0]
			} else if vals := md.Get("trace-id"); len(vals) > 0 {
				trace = vals[0]
			}
		}
		if trace == "" {
			trace = uuid.NewString()
		}
		ctx = errors.WithTraceID(ctx, trace)
		return handler(ctx, req)
	}
}

// httpStatusToGRPCCode maps HTTP-like status to gRPC code
func httpStatusToGRPCCode(httpStatus int) codes.Code {
	switch {
	case httpStatus == 401:
		return codes.Unauthenticated
	case httpStatus == 403:
		return codes.PermissionDenied
	case httpStatus >= 400 && httpStatus < 500:
		return codes.InvalidArgument
	default:
		return codes.Internal
	}
}

// UnaryServerInterceptor returns a gRPC interceptor with metrics + access logging + error translation.
// - serviceName: name for metrics/logging
// - store: domain.SessionStore (optional) used to enrich ctx with user_id and client IP
func UnaryServerInterceptor(serviceName string, store domain.SessionStore) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		fullMethod := info.FullMethod
		handlerName := handlerNameFromFullMethod(fullMethod)

		// metrics: inflight
		metrics.InflightRequests.WithLabelValues(serviceName, handlerName).Inc()
		defer metrics.InflightRequests.WithLabelValues(serviceName, handlerName).Dec()

		// 1) extract session-id from incoming metadata (if any)
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if svals := md.Get("x-session-id"); len(svals) > 0 {
				ctx = errors.WithSessionID(ctx, svals[0])
			} else if svals := md.Get("session-id"); len(svals) > 0 {
				ctx = errors.WithSessionID(ctx, svals[0])
			}
		}

		// 2) enrich ctx with user_id and optional clientIP from session store (best-effort)
		var sessionIP string
		if sid := errors.SessionIDFromContext(ctx); sid != "" && store != nil {
			if se, err := store.GetSession(ctx, sid); err == nil && se != nil {
				// store.GetSession returns domain.SessionEntry; keep only non-sensitive info
				if se.UserID != "" {
					ctx = errors.WithUserID(ctx, se.UserID)
				}
				if se.ClientIP != "" {
					sessionIP = se.ClientIP
				}
			}
		}

		// 3) get peer IP as fallback (best-effort)
		var peerIP string
		if p, ok := peer.FromContext(ctx); ok && p != nil {
			if addr, ok := p.Addr.(*net.TCPAddr); ok {
				peerIP = addr.IP.String()
			} else {
				// fallback to string (may include port)
				peerIP = p.Addr.String()
				if host, _, err := net.SplitHostPort(peerIP); err == nil {
					peerIP = host
				}
			}
		}
		// prefer session IP if available (set by store)
		if sessionIP != "" {
			peerIP = sessionIP
		}

		// 4) call actual handler
		resp, err := handler(ctx, req)
		latencyMs := time.Since(start).Milliseconds()

		// 5) error handling + metrics
		statusLabel := "OK"
		var appErr *errors.AppError
		if err != nil {
			appErr = errors.TranslateDomainError(ctx, err)
			if appErr == nil {
				appErr = &errors.AppError{
					Code:    "UNKNOWN",
					Status:  500,
					Message: "internal error",
					TraceID: errors.TraceIDFromContext(ctx),
				}
			}
			if st, ok := status.FromError(err); ok {
				statusLabel = st.Code().String()
			} else {
				statusLabel = "ERROR"
			}
			metrics.ErrorsTotal.WithLabelValues(serviceName, handlerName, appErr.Code).Inc()
		}

		metrics.RequestsTotal.WithLabelValues(serviceName, handlerName, fullMethod, statusLabel).Inc()
		metrics.RequestDuration.WithLabelValues(serviceName, handlerName, fullMethod).Observe(float64(latencyMs) / 1000.0)

		// 6) prepare access log fields (no duplication here)
		// env/instance/user/session/trace are added by logger.LogAccess by reading ctx and helpers,
		// so we avoid re-adding them here to prevent duplicates.
		httpCode := 200
		if err != nil {
			switch statusLabel {
			case "InvalidArgument":
				httpCode = 400
			case "Unauthenticated":
				httpCode = 401
			case "PermissionDenied":
				httpCode = 403
			case "NotFound":
				httpCode = 404
			default:
				httpCode = 500
			}
		}

		accessParams := logger.AccessParams{
			Service:   serviceName,
			Handler:   handlerName,
			Method:    fullMethod,
			HTTPCode:  httpCode,
			Status:    statusLabel,
			LatencyMS: latencyMs,
			IP:        peerIP,
			// keep Extra empty here to avoid duplicate instance/env/session/user entries;
			// use logger.LogAccess's internal enrichment.
			Extra: map[string]interface{}{},
		}

		logger.LogAccess(ctx, accessParams)

		// 7) translate domain app error -> grpc status for return
		if err == nil {
			return resp, nil
		}
		grpcCode := httpStatusToGRPCCode(appErr.Status)
		st := status.New(grpcCode, strings.Join([]string{appErr.Code, appErr.Message}, "|"))
		return nil, st.Err()
	}
}

// handlerNameFromFullMethod extracts short handler name
func handlerNameFromFullMethod(fullMethod string) string {
	if fullMethod == "" {
		return "unknown"
	}
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	parts := strings.SplitN(fullMethod, "/", 2)
	if len(parts) == 2 {
		servicePart := parts[0]
		if idx := strings.LastIndex(servicePart, "."); idx != -1 && idx+1 < len(servicePart) {
			servicePart = servicePart[idx+1:]
		}
		return servicePart + "/" + parts[1]
	}
	if idx := strings.LastIndex(fullMethod, "."); idx != -1 && idx+1 < len(fullMethod) {
		return fullMethod[idx+1:]
	}
	return fullMethod
}
