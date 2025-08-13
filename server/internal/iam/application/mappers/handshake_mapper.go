package mappers

import (
	"myapp/internal/iam/application/commands"
	"myapp/internal/iam/application/dtos"
)

func MapHandshakeRequestToCommand(dto dtos.HandshakeRequestDTO) commands.HandshakeCommand {
	return commands.HandshakeCommand{
		ClientPublicKey: dto.ClientPublicKey,
	}
}

func MapHandshakeResultToResponseDTO(result *commands.HandshakeResult) dtos.HandshakeResponseDTO {
	return dtos.HandshakeResponseDTO{
		ServerPublicKey:      result.ServerPublicKey,
		EncryptedSessionData: result.EncryptedSessionData,
		SessionID:            result.SessionID,
	}
}
