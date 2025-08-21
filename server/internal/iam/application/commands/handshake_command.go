package commands

type HandshakeCommand struct {
	ClientPublicKey string
}

type HandshakeResult struct {
	ServerPublicKey      string `json:"server_public_key"`
	SessionID            string `json:"session_id"`
	HKDFSaltB64          string `json:"hkdf_salt_b64"`
	ExpiresAt            int64  `json:"expires_at,omitempty"`
	EncryptedSessionData string `json:"encrypted_session_data,omitempty"`
}
