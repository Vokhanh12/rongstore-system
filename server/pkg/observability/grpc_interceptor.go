package observability

import (
	"context"
	"net"
	"os"
	"strings"
	"time"

	"server/internal/iam/domain"
	"server/pkg/auth"
	"server/pkg/logger"
	"server/pkg/metrics"
	"server/pkg/util/ctxutil"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var hostnameCache string

func init() {
	if hn, err := os.Hostname(); err == nil {
		hostnameCache = hn
	} else {
		hostnameCache = "unknown"
	}
}

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
		ctx = ctxutil.WithTraceID(ctx, trace)
		return handler(ctx, req)
	}
}

// Chuyển HTTP status sang gRPC code
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

func UnaryServerInterceptor(serviceName string, store domain.SessionStore, rules []auth.GrpcRule) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		fullMethod := info.FullMethod
		handlerName := handlerNameFromFullMethod(fullMethod)
		traceID := ctxutil.TraceIDFromContext(ctx)

		metrics.InflightRequests.WithLabelValues(serviceName, handlerName).Inc()
		defer metrics.InflightRequests.WithLabelValues(serviceName, handlerName).Dec()

		md, _ := metadata.FromIncomingContext(ctx)
		if err := auth.ValidateWithMetadata(md, fullMethod, rules); err != nil {
			logger.LogAccess(ctx, logger.AccessParams{
				Service:  serviceName,
				Handler:  handlerName,
				Method:   fullMethod,
				HTTPCode: 401,
				Status:   "Unauthenticated",
				Extra: map[string]interface{}{
					"reason": err.Error(),
				},
			})
			return nil, err
		}

		// if md, ok := metadata.FromIncomingContext(ctx); ok {
		// 	if err := ctxutil.ValidateMetadata(md, fullMethod, rules); err != nil {
		// 		logger.LogAccess(ctx, logger.AccessParams{
		// 			Service:  serviceName,
		// 			Handler:  handlerName,
		// 			Method:   fullMethod,
		// 			HTTPCode: 401,
		// 			Status:   "Unauthenticated",
		// 			Extra: map[string]interface{}{
		// 				"reason": err.Error(),
		// 			},
		// 		})
		// 		return nil, err
		// 	}

		// 	if svals := md.Get("x-session-id"); len(svals) > 0 {
		// 		ctx = ctxutil.WithSessionID(ctx, svals[0])
		// 	} else if svals := md.Get("session-id"); len(svals) > 0 {
		// 		ctx = ctxutil.WithSessionID(ctx, svals[0])
		// 	}
		// }

		if sid := ctxutil.SessionIDFromContext(ctx); sid != "" && store != nil {
			if se, err := store.GetSession(ctx, sid); err == nil && se != nil && se.UserID != "" {
				ctx = ctxutil.WithUserID(ctx, se.UserID)
			}
		}

		var peerIP string
		if p, ok := peer.FromContext(ctx); ok && p != nil {
			if addr, ok := p.Addr.(*net.TCPAddr); ok {
				peerIP = addr.IP.String()
			} else {
				peerIP = p.Addr.String()
				if host, _, err := net.SplitHostPort(peerIP); err == nil {
					peerIP = host
				}
			}
		}

		resp, err := handler(ctx, req)
		latencyMs := time.Since(start).Milliseconds()

		// httpCode := 200
		// statusLabel := "OK"
		// var errorCode string

		// if err != nil {
		// 	var be *ctxutil.BusinessError
		// 	if errors.As(err, &be) && be != nil {
		// 		httpCode = be.Status
		// 		statusLabel = "Error"
		// 		errorCode = be.Code
		// 		grpcCode := httpStatusToGRPCCode(be.Status)
		// 		st := status.New(grpcCode, fmt.Sprintf("%s|%s", be.Code, be.Message))
		// 		err = st.Err()
		// 		logger.LogError(ctx, handlerName, be.Message, "", be.WithExtra(map[string]interface{}{"method": fullMethod}).Data)
		// 	} else {
		// 		be = &ctxutil.BusinessError{
		// 			Code:      "CORE-INF-000",
		// 			Message:   err.Error(),
		// 			Status:    500,
		// 			Severity:  "S1",
		// 			Retryable: true,
		// 		}
		// 		httpCode = 500
		// 		statusLabel = "InternalError"
		// 		errorCode = be.Code
		// 		err = status.Errorf(codes.Internal, err.Error())
		// 		logger.LogError(ctx, handlerName, be.Message, "", be.WithExtra(map[string]interface{}{"method": fullMethod}).Data)
		// 	}
		// }

		// Metrics
		metrics.RequestsTotal.WithLabelValues(serviceName, handlerName, fullMethod, statusLabel).Inc()
		metrics.RequestDuration.WithLabelValues(serviceName, handlerName, fullMethod).Observe(float64(latencyMs) / 1000.0)

		// Logging access
		logger.LogAccess(ctx, logger.AccessParams{
			Service:   serviceName,
			Handler:   handlerName,
			Method:    fullMethod,
			HTTPCode:  httpCode,
			Status:    statusLabel,
			LatencyMS: latencyMs,
			IP:        peerIP,
			Extra: map[string]interface{}{
				"trace_id":   traceID,
				"error_code": errorCode,
			},
		})

		return resp, err
	}
}

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
