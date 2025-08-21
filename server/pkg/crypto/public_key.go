package crypto

import (
	"crypto/ecdh"
	"encoding/base64"
	"errors"
)

// ParsePublicKeyFromBase64Ptr decodes base64 and returns *ecdh.PublicKey (nilable)
func ParsePublicKeyFromBase64(curve ecdh.Curve, b64 string) (*ecdh.PublicKey, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	pubVal, err := curve.NewPublicKey(raw) // pubVal is ecdh.PublicKey (value)
	if err != nil {
		return nil, err
	}
	// take address of the value and return pointer
	return pubVal, nil
}

// Basic errors for domain compatibility (if domain.NewBusinessError expects more complex type, adjust accordingly)
var ErrBadPublicKey = errors.New("bad public key")
