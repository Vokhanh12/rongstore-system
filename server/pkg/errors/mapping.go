package errors

import (
	"net/http"
	dom "server/internal/iam/domain"
)

// domainToTransport là fallback tĩnh (nếu loader từ YAML không có key).
// Đảm bảo tên biến là domainToTransport (khớp với loader.go).
var domainToTransport = map[string]struct {
	Code    string
	Status  int
	Message string
}{
	dom.HandshakeInvalidClientKey: {Code: "AUTH-HAND-001", Status: http.StatusBadRequest, Message: "Invalid client public key"},
	dom.HandshakeKeyAgreementFail: {Code: "AUTH-HAND-002", Status: http.StatusUnprocessableEntity, Message: "Key agreement failed"},
	dom.HandshakeRNGFail:          {Code: "AUTH-HAND-003", Status: http.StatusInternalServerError, Message: "Server crypto failure"},
	dom.HandshakeEncryptFail:      {Code: "AUTH-HAND-004", Status: http.StatusInternalServerError, Message: "Failed to encrypt session"},

	// thêm mapping tĩnh khác nếu cần
}
