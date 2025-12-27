package grpc

import (
	"context"

	commonv1 "server/api/common/v1"
	iamv1 "server/api/iam/v1"
	rongstorev1 "server/api/rongstore/v1/resources"

	ucsiam "server/internal/iam/application/usecases"
	ucsrongstore "server/internal/iam/application/usecases"

	"server/pkg/logger"
	reshelper "server/pkg/util/response_helper"
)

type IamHandler struct {
	iamv1.UnimplementedIamServiceServer
	loginUsecase            *ucsiam.LoginUsecase
	handshakeUsecase        *ucsiam.HandshakeUsecase
	storeOwnerMutateUsecase *ucsiam.StoreOwnerMutateUsecase
}

func NewIamHandler(
	loginUsecase *ucsiam.LoginUsecase,
	handshakeUsecase *ucsiam.HandshakeUsecase,
	storeOwnerMutateUsecase *ucsiam.StoreOwnerMutateUsecase,
) *IamHandler {
	return &IamHandler{
		loginUsecase:            loginUsecase,
		handshakeUsecase:        handshakeUsecase,
		storeOwnerMutateUsecase: storeOwnerMutateUsecase,
	}
}

func (h *IamHandler) Login(ctx context.Context, req *iamv1.LoginRequest) (*commonv1.BaseResponse, error) {
	cmd := ucsiam.MapLoginRequestToCommand(req)

	result, err := h.loginUsecase.Execute(ctx, cmd)
	if err != nil {

		logger.LogBySeverity(ctx, "iam_handler.login", err)

		return reshelper.BuildErrorResponse(ctx, err), nil
	}

	resDTO := ucsiam.MapLoginResultToResponseDTO(result)
	return reshelper.BuildSuccessResponse(ctx, &resDTO)
}

func (h *IamHandler) Handshake(ctx context.Context, req *iamv1.HandshakeRequest) (*commonv1.BaseResponse, error) {
	cmd := ucsiam.MapHandshakeRequestToCommand(req)

	result, err := h.handshakeUsecase.Execute(ctx, cmd)
	if err != nil {

		logger.LogBySeverity(ctx, "iam_handler.handshake", err)

		return reshelper.BuildErrorResponse(ctx, err), nil
	}

	resDTO := ucsiam.MapHandshakeResultToResponseDTO(result)
	return reshelper.BuildSuccessResponse(ctx, &resDTO)
}

func (h *IamHandler) StoreOwnerMutate(ctx context.Context, req *rongstorev1.StoreOwnerMutateRequest) (*commonv1.MutateResponse, error) {

	cmd := ucsrongstore.MapStoreOwnerMutateRequestToCommand(req)

	result, err := h.storeOwnerMutateUsecase.Execute(ctx, cmd)
	if err != nil {
		logger.LogBySeverity(ctx, "iam_handler.store_owner_mutate", err)
		return reshelper.BuildErrorResponse(ctx, err), nil
	}

	resDTO := ucsrongstore.MapStoreOwnerMutateResultToResponseDTO(result)

	return &resDTO, nil
}
