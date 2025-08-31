package auth

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GrpcRule struct {
	MethodPrefix string
	Required     []string
}

func DefaultGrpcRules() []GrpcRule {
	return []GrpcRule{
		{MethodPrefix: "/iam.v1.IAMService/Handshake", Required: []string{"session-id", "x-session-id"}},
		{MethodPrefix: "*", Required: []string{"authorization"}},
	}
}

func ValidateWithMetadata(md metadata.MD, fullMethod string, rules []GrpcRule) error {
	lowerMD := map[string][]string{}
	for k, vals := range md {
		lowerMD[strings.ToLower(k)] = vals
	}

	var matched *GrpcRule
	for _, r := range rules {
		if r.MethodPrefix == "*" || strings.HasPrefix(fullMethod, r.MethodPrefix) {
			matched = &r
			break
		}
	}
	if matched == nil {
		return nil
	}

	for _, req := range matched.Required {
		if vals, ok := lowerMD[strings.ToLower(req)]; !ok || len(vals) == 0 || strings.TrimSpace(vals[0]) == "" {
			return status.Errorf(codes.Unauthenticated, "missing metadata: %s", req)
		}
	}
	return nil
}
