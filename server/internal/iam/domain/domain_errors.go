package domain

import "fmt"

// BusinessError đại diện cho lỗi ở tầng domain (business).
// - Key: machine-readable key (ví dụ "HANDSHAKE_INVALID_CLIENT_KEY")
// - Message: canonical message (có thể dùng để fallback cho client nếu không map được)
// - Data: optional map chứa metadata an toàn (KHÔNG chứa secret)
type BusinessError struct {
	Key     string
	Message string
	Data    map[string]interface{}
}

func (e *BusinessError) Error() string {
	if e == nil {
		return "<nil BusinessError>"
	}
	return fmt.Sprintf("%s: %s", e.Key, e.Message)
}

func NewBusinessError(key, message string, data map[string]interface{}) *BusinessError {
	return &BusinessError{
		Key:     key,
		Message: message,
		Data:    data,
	}
}

func (e *BusinessError) MessageOr(fallback string) string {
	if e == nil {
		return fallback
	}
	if e.Message == "" {
		return fallback
	}
	return e.Message
}

// --- Domain keys (const) ---
// Thêm các key ở đây; dùng chúng trong usecase để trả lỗi domain.
// Nếu cần thêm lỗi khác, bổ sung vào danh sách này.
const (
	// Handshake errors
	HandshakeInvalidClientKey = "HANDSHAKE_INVALID_CLIENT_KEY"
	HandshakeKeyAgreementFail = "HANDSHAKE_KEY_AGREEMENT_FAIL"
	HandshakeRNGFail          = "HANDSHAKE_RNG_FAIL"
	HandshakeEncryptFail      = "HANDSHAKE_ENCRYPT_FAIL"

	// User / Auth examples (bạn có thể mở rộng)
	UserNotFound   = "USER_NOT_FOUND"
	UserEmailTaken = "USER_EMAIL_TAKEN"
)
