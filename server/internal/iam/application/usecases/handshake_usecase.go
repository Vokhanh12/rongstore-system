package usecases

import (
	"context"
	"crypto/aes"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"time"

	"server/internal/iam/application/commands"
	"server/internal/iam/domain"
	"server/pkg/crypto"
)

// zeroBytes zeroes the provided byte slice in-place.
func zeroBytes(b []byte) {
	if b == nil {
		return
	}
	for i := range b {
		b[i] = 0
	}
}

type HandshakeUsecase struct {
	UserRepo domain.UserRepository
}

func NewHandshakeUsecase(repo domain.UserRepository) *HandshakeUsecase {
	return &HandshakeUsecase{
		UserRepo: repo,
	}
}

func (u *HandshakeUsecase) Execute(ctx context.Context, cmd commands.HandshakeCommand) (*commands.HandshakeResult, error) {
	// Basic validation
	if cmd.ClientPublicKey == "" {
		return nil, domain.NewBusinessError(domain.HandshakeInvalidClientKey, "client public key is required", nil)
	}

	curve := ecdh.P521()

	// 🔐 Tạo key pair server
	serverPriv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		// server RNG/crypto failure
		return nil, domain.NewBusinessError(domain.HandshakeRNGFail, "server key generation failed", map[string]interface{}{"cause": err.Error()})
	}

	// ensure server private key bytes are not leaked (serverPriv is a crypto.PrivateKey interface; zeroing handled by GC)
	// 📥 Parse client public key
	clientPub, err := crypto.ParsePublicKeyFromBase64(cmd.ClientPublicKey)
	if err != nil {
		// client-supplied invalid key
		return nil, domain.NewBusinessError(domain.HandshakeInvalidClientKey, "invalid client public key", map[string]interface{}{"client_pub_len": len(cmd.ClientPublicKey)})
	}

	// 🔐 Derive shared secret
	sharedSecret, err := serverPriv.ECDH(clientPub)
	if err != nil {
		// key agreement failed (client key may be incompatible)
		return nil, domain.NewBusinessError(domain.HandshakeKeyAgreementFail, "key agreement failed", nil)
	}
	// Ensure we wipe sharedSecret when done
	defer zeroBytes(sharedSecret)

	// 🔐 Sinh AES key, IV, salt
	aesKey := make([]byte, 32) // AES-256
	iv := make([]byte, aes.BlockSize)
	salt := make([]byte, 16)

	if _, err := rand.Read(aesKey); err != nil {
		zeroBytes(aesKey)
		return nil, domain.NewBusinessError(domain.HandshakeRNGFail, "random generation failed for AES key", map[string]interface{}{"cause": err.Error()})
	}
	if _, err := rand.Read(iv); err != nil {
		zeroBytes(aesKey)
		zeroBytes(iv)
		return nil, domain.NewBusinessError(domain.HandshakeRNGFail, "random generation failed for IV", map[string]interface{}{"cause": err.Error()})
	}
	if _, err := rand.Read(salt); err != nil {
		zeroBytes(aesKey)
		zeroBytes(iv)
		zeroBytes(salt)
		return nil, domain.NewBusinessError(domain.HandshakeRNGFail, "random generation failed for salt", map[string]interface{}{"cause": err.Error()})
	}

	// 🧱 Chuẩn bị session data để mã hóa
	sessionInfo := crypto.SessionInfo{
		AESKey:     aesKey,
		IV:         iv,
		Salt:       salt,
		Expiration: time.Now().Add(15 * time.Minute),
	}

	// 🔐 Mã hóa sessionInfo bằng sharedSecret
	encryptedData, err := crypto.EncryptSessionInfo(sharedSecret, sessionInfo)
	// wipe AES buffers once encrypted (we still keep encryptedData)
	zeroBytes(aesKey)
	zeroBytes(iv)
	zeroBytes(salt)
	if err != nil {
		return nil, domain.NewBusinessError(domain.HandshakeEncryptFail, "failed to encrypt session info", map[string]interface{}{"cause": err.Error()})
	}

	// 🧾 Tạo session ID (simple UUID string)
	sessionID := crypto.GenerateSessionID()

	// 🧬 Encode public key
	serverPubKey := serverPriv.PublicKey()
	serverPubKeyBytes := serverPubKey.Bytes()
	serverPubKeyBase64 := base64.StdEncoding.EncodeToString(serverPubKeyBytes)

	// NOTE: Do NOT include any secret (sharedSecret or raw AES key) in the response or error data.
	// sharedSecret will be zeroed by defer above.

	return &commands.HandshakeResult{
		ServerPublicKey:      serverPubKeyBase64,
		EncryptedSessionData: base64.StdEncoding.EncodeToString(encryptedData),
		SessionID:            sessionID,
	}, nil
}
