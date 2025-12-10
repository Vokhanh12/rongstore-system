package logger

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func MaskEmail(email string) string {
	if !strings.Contains(email, "@") {
		return email
	}
	parts := strings.Split(email, "@")
	return parts[0][:1] + "***@" + parts[1]
}

func TokenHashPrefix(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])[:10]
}
