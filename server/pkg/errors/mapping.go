package errors

import (
	"net/http"
	dom "server/internal/iam/domain"
)

var handshakeMap = map[string]struct {
	Code    string
	Status  int
	Message string
}{
	dom.HandshakeInvalidClientKey: {Code: "AUTH-HAND-001", Status: http.StatusBadRequest, Message: "Invalid client public key"},
	dom.HandshakeKeyAgreementFail: {Code: "AUTH-HAND-002", Status: http.StatusBadRequest, Message: "Key agreement failed"},
	dom.HandshakeRNGFail:          {Code: "AUTH-HAND-003", Status: http.StatusInternalServerError, Message: "Server crypto failure"},
	dom.HandshakeEncryptFail:      {Code: "AUTH-HAND-004", Status: http.StatusInternalServerError, Message: "Failed to encrypt session"},
}
