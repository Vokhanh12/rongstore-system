package mappers

import (
	iamv1 "server/api/iam/v1"
	"server/internal/iam/application/commands"
)

func MapHandshakeRequestToCommand(req *iamv1.HandshakeRequest) commands.HandshakeCommand {
	return commands.HandshakeCommand{
		ClientPublicKey: req.ClientPublicKey,
	}
}

func MapHandshakeResultToResponseDTO(result *commands.HandshakeResult) iamv1.HandshakeResponse {
	return iamv1.HandshakeResponse{
		ServerPublicKey:      result.ServerPublicKey,
		SessionId:            result.SessionID,
		HkdfSaltB64:          result.HKDFSaltB64,
		ExpiresAt:            int32(result.ExpiresAt),
		EncryptedSessionData: result.EncryptedSessionData,
	}
}
