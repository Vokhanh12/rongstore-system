package usecases

import (
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
)

// zeroBytes overwrites sensitive data in memory.
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
	// basic validation
	if cmd.ClientPublicKey == "" {
		return nil, domain.NewBusinessError(domain.HandshakeInvalidClientKey, "client public key is required", nil)
	}

	// Choose curve (P-521)
	curve := ecdh.P521()

	// 1) server ephemeral keypair
	serverPriv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, domain.NewBusinessError(domain.HandshakeRNGFail, "server key generation failed", map[string]interface{}{"cause": err.Error()})
	}
	serverPub := serverPriv.PublicKey()
	serverPubBytes := serverPub.Bytes()
	serverPubB64 := base64.StdEncoding.EncodeToString(serverPubBytes)

	// 2) parse client public key
	clientPub, err := crypto.ParsePublicKeyFromBase64(curve, cmd.ClientPublicKey)
	if err != nil {
		return nil, domain.NewBusinessError(
			domain.HandshakeInvalidClientKey,
			"invalid client public key",
			map[string]interface{}{"client_pub_len": len(cmd.ClientPublicKey)},
		)
	}
	clientPubBytes := clientPub.Bytes()

	// 3) ECDH shared secret
	sharedSecret, err := serverPriv.ECDH(clientPub)
	if err != nil {
		return nil, domain.NewBusinessError(domain.HandshakeKeyAgreementFail, "key agreement failed", nil)
	}
	defer zeroBytes(sharedSecret)

	// 4) hkdf salt + sessionID
	hkdfSalt := make([]byte, 16)
	if _, err := rand.Read(hkdfSalt); err != nil {
		return nil, domain.NewBusinessError(domain.HandshakeRNGFail, "random generation failed for hkdf salt", map[string]interface{}{"cause": err.Error()})
	}
	sessionID := uuid.NewString()
	expiry := time.Now().Add(15 * time.Minute)

	// 5) info context
	info := make([]byte, 0, 256)
	info = append(info, []byte("handshake|derive|v1|")...)
	info = append(info, serverPubBytes...)
	info = append(info, byte('|'))
	info = append(info, clientPubBytes...)
	info = append(info, byte('|'))
	info = append(info, []byte(sessionID)...)

	// 6) derive kc2s and ks2c
	hk := hkdf.New(sha256.New, sharedSecret, hkdfSalt, append(info, []byte("|c2s")...))
	kc2s := make([]byte, 32)
	if _, err := io.ReadFull(hk, kc2s); err != nil {
		zeroBytes(kc2s)
		return nil, domain.NewBusinessError(domain.HandshakeKeyDeriveFail, "hkdf derive failed (c2s)", map[string]interface{}{"cause": err.Error()})
	}

	hk2 := hkdf.New(sha256.New, sharedSecret, hkdfSalt, append(info, []byte("|s2c")...))
	ks2c := make([]byte, 32)
	if _, err := io.ReadFull(hk2, ks2c); err != nil {
		zeroBytes(kc2s)
		zeroBytes(ks2c)
		return nil, domain.NewBusinessError(domain.HandshakeKeyDeriveFail, "hkdf derive failed (s2c)", map[string]interface{}{"cause": err.Error()})
	}

	// 7) store session via SessionStore (Redis or memory)
	entry := &domain.SessionEntry{
		SessionID: sessionID,
		ClientPub: clientPubBytes,
		ServerPub: serverPubBytes,
		Kc2s:      kc2s,
		Ks2c:      ks2c,
		HKDFSalt:  hkdfSalt,
		Expiry:    expiry,
	}
	if err := u.SessionStore.StoreSession(ctx, entry); err != nil {
		return nil, domain.NewBusinessError(domain.HandshakeStorageFail, "failed to store session", map[string]interface{}{"cause": err.Error()})
	}

	// 8) build response
	return &commands.HandshakeResult{
		ServerPublicKey: serverPubB64,
		SessionID:       sessionID,
		HKDFSaltB64:     base64.StdEncoding.EncodeToString(hkdfSalt),
		ExpiresAt:       expiry.Unix(),
	}, nil
}
