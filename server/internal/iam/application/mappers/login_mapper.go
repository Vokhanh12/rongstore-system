package mappers

import (
	iamv1 "server/api/iam/v1"
	"server/internal/iam/application/commands"
)

func MapLoginRequestToCommand(req *iamv1.LoginRequest) commands.LoginCommand {
	return commands.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	}
}

func MapLoginResultToResponseDTO(result *commands.LoginResult) iamv1.LoginResponse {
	return iamv1.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}
}
