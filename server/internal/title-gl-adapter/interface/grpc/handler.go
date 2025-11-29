package grpc

import (
	"context"

	iamv1 "server/api/iam/v1"
	"server/internal/iam/application/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IamHandler struct {
	iamv1.UnimplementedIamServiceServer
	loginUsecase     *usecases.LoginUserUsecase
	handshakeUsecase *usecases.HandshakeUsecase
}

func NewIamHandler(
	loginUsecase *usecases.LoginUserUsecase,
	handshakeUsecase *usecases.HandshakeUsecase,
) *IamHandler {
	return &IamHandler{
		loginUsecase:     loginUsecase,
		handshakeUsecase: handshakeUsecase,
	}
}

func (h *IamHandler) Login(ctx context.Context, req *iamv1.LoginRequest) (*iamv1.LoginResponse, error) {

	cmd := usecases.MapLoginRequestToCommand(req)

	result, err := h.loginUsecase.Execute(ctx, cmd)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	resDTO := usecases.MapLoginResultToResponseDTO(result)
	return &iamv1.LoginResponse{
		AccessToken:  resDTO.AccessToken,
		RefreshToken: resDTO.RefreshToken,
	}, nil
}

func (h *IamHandler) Handshake(ctx context.Context, req *iamv1.HandshakeRequest) (*iamv1.HandshakeResponse, error) {

	cmd := usecases.MapHandshakeRequestToCommand(req)

	result, err := h.handshakeUsecase.Execute(ctx, cmd)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resDTO := usecases.MapHandshakeResultToResponseDTO(result)
	return &iamv1.HandshakeResponse{
		ServerPublicKey:      resDTO.ServerPublicKey,
		SessionId:            resDTO.SessionId,
		HkdfSaltB64:          resDTO.HkdfSaltB64,
		ExpiresAt:            resDTO.ExpiresAt,
		EncryptedSessionData: resDTO.EncryptedSessionData,
	}, nil
}
