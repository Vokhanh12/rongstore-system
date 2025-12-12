package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc/peer"
)

func simplifyMethod(fullMethod string) string {
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	parts := strings.SplitN(fullMethod, "/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return fullMethod
}

func peerIP(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok && p != nil {
		return p.Addr.String()
	}
	return ""
}
