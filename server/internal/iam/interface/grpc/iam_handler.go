package grpc

import (
	"context"

	commonv1 "server/api/common/v1"
	iamv1 "server/api/iam/v1"

	usecases "server/internal/iam/application/usecases/auth"

	"server/pkg/util/grpcutil"
	reshelper "server/pkg/util/response_helper"
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
		return grpcutil.HandleBusinessError(ctx, "Login", req, err), nil
	}

	resDTO := usecases.MapLoginResultToResponseDTO(result)
	return reshelper.BuildSuccessResponse(ctx, &resDTO)
}

func (h *IamHandler) Handshake(ctx context.Context, req *iamv1.HandshakeRequest) (*commonv1.BaseResponse, error) {
	cmd := usecases.MapHandshakeRequestToCommand(req)

	result, err := h.handshakeUsecase.Execute(ctx, cmd)
	if err != nil {
		return grpcutil.HandleBusinessError(ctx, "Handshake", req, err), nil
	}

	resDTO := usecases.MapHandshakeResultToResponseDTO(result)
	return reshelper.BuildSuccessResponse(ctx, &resDTO)
}
