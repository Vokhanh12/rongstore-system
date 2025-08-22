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
		return nil, domain.NewBusinessError(domain.HandshakeInvalidClientKey, "client public key is required", nil)
	}

	// 1) Curve: X25519 (nhanh + phổ biến)
	curve := ecdh.X25519()

	// 2) Server ephemeral keypair
	serverPriv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, domain.NewBusinessError(domain.HandshakeRNGFail, "server key generation failed", map[string]interface{}{"cause": err.Error()})
	}
	serverPub := serverPriv.PublicKey()
	serverPubBytes := serverPub.Bytes() // 32B for X25519
	serverPubB64 := base64.StdEncoding.EncodeToString(serverPubBytes)

	// 3) Parse client public key
	clientPub, err := crypto.ParsePublicKeyFromBase64(curve, cmd.ClientPublicKey)
	if err != nil {
		return nil, domain.NewBusinessError(
			domain.HandshakeInvalidClientKey,
			"invalid client public key",
			map[string]interface{}{"client_pub_len": len(cmd.ClientPublicKey)},
		)
	}
	clientPubBytes := clientPub.Bytes()

	// 4) ECDH shared secret
	sharedSecret, err := serverPriv.ECDH(clientPub)
	if err != nil {
		return nil, domain.NewBusinessError(domain.HandshakeKeyAgreementFail, "key agreement failed", nil)
	}
	// Xóa ngay khi dùng xong (best-effort)
	defer zeroBytes(sharedSecret)

	// 5) HKDF salt (32 bytes = HashLen của SHA-256)
	hkdfSalt := make([]byte, 32)
	if _, err := rand.Read(hkdfSalt); err != nil {
		return nil, domain.NewBusinessError(domain.HandshakeRNGFail, "random generation failed for hkdf salt", map[string]interface{}{"cause": err.Error()})
	}

	// 6) Session metadata
	sessionID := uuid.NewString()
	expiry := time.Now().UTC().Add(15 * time.Minute)

	// 7) info: gắn version + suite để tách bạch ngữ cảnh lâu dài
	var infoBuf bytes.Buffer
	infoBuf.WriteString("handshake|derive|v1|aead=aes-gcm|hash=sha256|curve=x25519|")
	infoBuf.Write(serverPubBytes)
	infoBuf.WriteByte('|')
	infoBuf.Write(clientPubBytes)
	infoBuf.WriteByte('|')
	infoBuf.WriteString(sessionID)
	info := infoBuf.Bytes()

	// 8) Derive một lượt: 32+32+12+12 = 88 bytes
	// kc2s, ks2c: 32B (AES-256-GCM hoặc ChaCha20-Poly1305)
	// nonceC2S, nonceS2C: 12B (GCM/ChaCha20-Poly1305 nonce)
	okm := make([]byte, 88)
	hk := hkdf.New(sha256.New, sharedSecret, hkdfSalt, info)
	if _, err := io.ReadFull(hk, okm); err != nil {
		zeroBytes(okm)
		return nil, domain.NewBusinessError(domain.HandshakeKeyDeriveFail, "hkdf derive failed", map[string]interface{}{"cause": err.Error()})
	}

	kc2s := okm[0:32]
	ks2c := okm[32:64]
	// nonceC2S := okm[64:76]
	// nonceS2C := okm[76:88]

	// 9) Store session (Redis/memory)
	entry := &domain.SessionEntry{
		SessionID: sessionID,
		ClientPub: clientPubBytes,
		ServerPub: serverPubBytes,
		Kc2s:      append([]byte(nil), kc2s...), // clone để tránh alias
		Ks2c:      append([]byte(nil), ks2c...),
		HKDFSalt:  append([]byte(nil), hkdfSalt...),
		Expiry:    expiry,

		// Khuyến nghị: thêm 2 trường này vào SessionEntry để dùng nonce-base per-direction
		// NonceC2S:   append([]byte(nil), nonceC2S...),
		// NonceS2C:   append([]byte(nil), nonceS2C...),
	}
	if err := u.SessionStore.StoreSession(ctx, entry); err != nil {
		zeroBytes(okm)
		zeroBytes(kc2s)
		zeroBytes(ks2c)
		return nil, domain.NewBusinessError(domain.HandshakeStorageFail, "failed to store session", map[string]interface{}{"cause": err.Error()})
	}

	// 10) Xoá tạm bộ nhớ nhạy cảm
	zeroBytes(okm) // đã clone các phần cần lưu
	// kc2s/ks2c trong entry cần giữ để dùng cho phiên; bản tạm đã xoá okm ở trên

	// 11) Response
	return &commands.HandshakeResult{
		ServerPublicKey: serverPubB64,
		SessionID:       sessionID,
		HKDFSaltB64:     base64.StdEncoding.EncodeToString(hkdfSalt),
		ExpiresAt:       expiry.Unix(), // seconds, UTC
	}, nil
}
