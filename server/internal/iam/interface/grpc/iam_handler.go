package grpc

import (
	"context"

	commonv1 "server/api/common/v1"
	iamv1 "server/api/iam/v1"

	usecases "server/internal/iam/application/usecases"
	"server/internal/iam/domain/services"

	"server/pkg/logger"
	reshelper "server/pkg/util/response_helper"
)

type IamHandler struct {
	iamv1.UnimplementedIamServiceServer
	business_errors  services.BusinessError
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

		businessError := h.business_errors.GetBusinessError(err)
		logger.LogBySeverity(ctx, *businessError, map[string]interface{}{
			"handler": "Login",
			"request": req,
		})

		return reshelper.BuildErrorResponse(ctx, businessError), nil
	}

	resDTO := usecases.MapLoginResultToResponseDTO(result)
	return reshelper.BuildSuccessResponse(ctx, &resDTO)
}

func (h *IamHandler) Handshake(ctx context.Context, req *iamv1.HandshakeRequest) (*commonv1.BaseResponse, error) {
	cmd := usecases.MapHandshakeRequestToCommand(req)

	result, err := h.handshakeUsecase.Execute(ctx, cmd)
	if err != nil {
		businessError := h.business_errors.GetBusinessError(err)
		logger.LogBySeverity(ctx, *businessError, map[string]interface{}{
			"handler": "Handshake",
			"request": req,
		})

		return reshelper.BuildErrorResponse(ctx, businessError), nil
	}

	resDTO := usecases.MapHandshakeResultToResponseDTO(result)
	return reshelper.BuildSuccessResponse(ctx, &resDTO)
}
