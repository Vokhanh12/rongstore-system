package grpc

import (
	"context"
	"time"

	commonv1 "server/api/common/v1"
	iamv1 "server/api/iam/v1"

	usecases "server/internal/iam/application/usecases/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

type IamHandler struct {
	iamv1.UnimplementedIamServiceServer
	loginUsecase     *usecases.LoginUsecase
	handshakeUsecase *usecases.HandshakeUsecase
}

func NewIamHandler(
	loginUsecase *usecases.LoginUsecase,
	handshakeUsecase *usecases.HandshakeUsecase,
) *IamHandler {
	return &IamHandler{
		loginUsecase:     loginUsecase,
		handshakeUsecase: handshakeUsecase,
	}
}

func (h *IamHandler) Login(ctx context.Context, req *iamv1.LoginRequest) (*commonv1.BaseResponse, error) {

	cmd := usecases.MapLoginRequestToCommand(req)

	result, err := h.loginUsecase.Execute(ctx, cmd)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	resDTO := usecases.MapLoginResultToResponseDTO(result)
	anyData, err := anypb.New(&resDTO)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal response")
	}

	return &commonv1.BaseResponse{
		Success: true,
		Data:    anyData,
		Metadata: &commonv1.Metadata{
			TraceId:   "trace-id-example",
			RequestId: "req-id-example",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}, nil

}

func (h *IamHandler) Handshake(ctx context.Context, req *iamv1.HandshakeRequest) (*commonv1.BaseResponse, error) {

	cmd := usecases.MapHandshakeRequestToCommand(req)

	result, err := h.handshakeUsecase.Execute(ctx, cmd)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resDTO := usecases.MapHandshakeResultToResponseDTO(result)
	anyData, err := anypb.New(&resDTO)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal response")
	}
	return &commonv1.BaseResponse{
		Success: true,
		Data:    anyData,
		Metadata: &commonv1.Metadata{
			TraceId:   "trace-id-example",
			RequestId: "req-id-example",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}, nil
}
