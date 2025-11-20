// Package response_helper provides utilities to build standard API responses.
package response_helper

import (
	"context"
	"time"

	"server/internal/hr/domain"
	domain_error "server/internal/iam/domain"

	"server/pkg/errors"
	"server/pkg/util/ctxutil"

	commonv1 "server/api/common/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// BuildErrorResponse builds a BaseResponse representing a failed operation.
// It sets success=false, fills Metadata, and sets Error code.
func BuildErrorResponse(ctx context.Context, be *errors.BusinessError) *commonv1.BaseResponse {
	if be == nil {
		be = &domain_error.UNKNOWN_DOMAIN_KEY
	}

	return &commonv1.BaseResponse{
		Success: false,
		Metadata: &commonv1.Metadata{
			TraceId:   ctxutil.TraceIDFromContext(ctx),
			RequestId: ctxutil.RequestIdFromContext(ctx),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Error: &commonv1.Error{
			Code:     be.Code,
			Message:  be.Message,
			Severity: be.Severity,
		},
	}
}

// BuildSuccessResponse builds a BaseResponse representing a successful operation.
// It converts the provided protobuf message to google.protobuf.Any.
func BuildSuccessResponse(ctx context.Context, data proto.Message) (*commonv1.BaseResponse, error) {
	anyData, err := anypb.New(data)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal response")
	}

	return &commonv1.BaseResponse{
		Success: true,
		Data:    anyData,
		Metadata: &commonv1.Metadata{
			TraceId:   ctxutil.TraceIDFromContext(ctx),
			RequestId: ctxutil.RequestIdFromContext(ctx),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}, nil
}

// FromError converts a generic error into a BaseResponse using BusinessError mapping.
func FromError(ctx context.Context, err error) *commonv1.BaseResponse {
	be, _ := domain.GetBusinessError(err)
	return BuildErrorResponse(ctx, be)
}
