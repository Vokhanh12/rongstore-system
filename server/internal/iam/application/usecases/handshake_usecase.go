package usecases

import (
	"bytes"
	"context"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/hkdf"

	"server/internal/iam/application/commands"
	"server/internal/iam/domain"
	"server/pkg/crypto"
	"server/pkg/errors"
)

// zeroBytes overwrites sensitive data in memory (best-effort in Go).
func zeroBytes(b []byte) {
	if b == nil {
		return
	}
	for i := range b {
		b[i] = 0
	}
}

type HandshakeUsecase struct {
	UserRepo     domain.UserRepository
	SessionStore domain.SessionStore
}

func NewHandshakeUsecase(repo domain.UserRepository, store domain.SessionStore) *HandshakeUsecase {
	return &HandshakeUsecase{
		UserRepo:     repo,
		SessionStore: store,
	}
}

func (u *HandshakeUsecase) Execute(ctx context.Context, cmd commands.HandshakeCommand) (*commands.HandshakeResult, error) {
	// 0) basic validation
	if cmd.ClientPublicKey == "" {
		return nil, errors.NewBusinessError(
			domain.HANDSHAKE_INVALID_CLIENT_KEY,
			errors.WithMessage("client public key is required"),
		)
	}

	// 1) Curve: X25519 (nhanh + phổ biến)
	curve := ecdh.X25519()

	// 2) Server ephemeral keypair
	serverPriv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, errors.NewBusinessError(
			domain.HANDSHAKE_RNG_FAIL,
			errors.WithMessage("server key generation failed"),
			errors.WithData(map[string]interface{}{"cause": err.Error()}),
		)
	}
	serverPub := serverPriv.PublicKey()
	serverPubBytes := serverPub.Bytes()
	serverPubB64 := base64.StdEncoding.EncodeToString(serverPubBytes)

	// 3) Parse client public key
	clientPub, err := crypto.ParsePublicKeyFromBase64(curve, cmd.ClientPublicKey)
	if err != nil {
		return nil, errors.NewBusinessError(
			domain.HANDSHAKE_INVALID_CLIENT_KEY,
			errors.WithMessage("invalid client public key"),
			errors.WithData(map[string]interface{}{"client_pub_len": len(cmd.ClientPublicKey)}),
		)
	}
	clientPubBytes := clientPub.Bytes()

	// 4) ECDH shared secret
	sharedSecret, err := serverPriv.ECDH(clientPub)
	if err != nil {
		return nil, errors.NewBusinessError(
			domain.HANDSHAKE_KEY_AGREEMENT_FAIL,
			errors.WithMessage("key agreement failed"),
		)
	}
	defer zeroBytes(sharedSecret)

	// 5) HKDF salt
	hkdfSalt := make([]byte, 32)
	if _, err := rand.Read(hkdfSalt); err != nil {
		return nil, errors.NewBusinessError(
			domain.HANDSHAKE_RNG_FAIL,
			errors.WithMessage("random generation failed for hkdf salt"),
			errors.WithData(map[string]interface{}{"cause": err.Error()}),
		)
	}

	// 6) Session metadata
	sessionID := uuid.NewString()
	expiry := time.Now().UTC().Add(15 * time.Minute)

	// 7) info: gắn version + suite
	var infoBuf bytes.Buffer
	infoBuf.WriteString("handshake|derive|v1|aead=aes-gcm|hash=sha256|curve=x25519|")
	infoBuf.Write(serverPubBytes)
	infoBuf.WriteByte('|')
	infoBuf.Write(clientPubBytes)
	infoBuf.WriteByte('|')
	infoBuf.WriteString(sessionID)
	info := infoBuf.Bytes()

	// 8) Derive keys
	okm := make([]byte, 88)
	hk := hkdf.New(sha256.New, sharedSecret, hkdfSalt, info)
	if _, err := io.ReadFull(hk, okm); err != nil {
		zeroBytes(okm)
		return nil, errors.NewBusinessError(
			domain.HANDSHAKE_KEY_DERIVE_FAIL,
			errors.WithMessage("hkdf derive failed"),
			errors.WithData(map[string]interface{}{"cause": err.Error()}),
		)
	}

	kc2s := okm[0:32]
	ks2c := okm[32:64]

	// 9) Store session
	entry := &domain.SessionEntry{
		SessionID: sessionID,
		ClientPub: clientPubBytes,
		ServerPub: serverPubBytes,
		Kc2s:      append([]byte(nil), kc2s...),
		Ks2c:      append([]byte(nil), ks2c...),
		HKDFSalt:  append([]byte(nil), hkdfSalt...),
		Expiry:    expiry,
	}
	if err := u.SessionStore.StoreSession(ctx, entry); err != nil {
		zeroBytes(okm)
		zeroBytes(kc2s)
		zeroBytes(ks2c)
		return nil, errors.NewBusinessError(
			domain.HANDSHAKE_STORAGE_FAIL,
			errors.WithMessage("failed to store session"),
			errors.WithData(map[string]interface{}{"cause": err.Error()}),
		)
	}

	zeroBytes(okm)

	return &commands.HandshakeResult{
		ServerPublicKey: serverPubB64,
		SessionID:       sessionID,
		HKDFSaltB64:     base64.StdEncoding.EncodeToString(hkdfSalt),
		ExpiresAt:       expiry.Unix(),
	}, nil
}
