package grpc

import (
	"context"

	commonv1 "server/api/common/v1"
	iamv1 "server/api/iam/v1"

	usecases "server/internal/iam/application/usecases"

	"server/pkg/logger"
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

		logger.LogBySeverity(ctx, "iam_handler.login", err, map[string]interface{}{
			"handler": "Login",
			"request": req,
		})

		return reshelper.BuildErrorResponse(ctx, err), nil
	}

	resDTO := usecases.MapLoginResultToResponseDTO(result)
	return reshelper.BuildSuccessResponse(ctx, &resDTO)
}

func (h *IamHandler) Handshake(ctx context.Context, req *iamv1.HandshakeRequest) (*commonv1.BaseResponse, error) {
	cmd := usecases.MapHandshakeRequestToCommand(req)

	result, err := h.handshakeUsecase.Execute(ctx, cmd)
	if err != nil {

		logger.LogBySeverity(ctx, "iam_handler.handshake", err, map[string]interface{}{
			"handler": "Handshake",
			"request": req,
		})

		return reshelper.BuildErrorResponse(ctx, err), nil
	}

	resDTO := usecases.MapHandshakeResultToResponseDTO(result)
	return reshelper.BuildSuccessResponse(ctx, &resDTO)
}
