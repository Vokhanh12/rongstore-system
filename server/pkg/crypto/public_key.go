package crypto

import (
	"crypto/ecdh"
	"encoding/base64"
	"errors"
)

func ParsePublicKeyFromBase64(curve ecdh.Curve, b64 string) (*ecdh.PublicKey, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	pubVal, err := curve.NewPublicKey(raw)
	if err != nil {
		return nil, err
	}
	return pubVal, nil
}

var ErrBadPublicKey = errors.New("bad public key")
