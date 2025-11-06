package crypto

// ZeroBytes overwrites sensitive data in memory (best-effort in Go)
func ZeroBytes(b []byte) {
	if b == nil {
		return
	}
	for i := range b {
		b[i] = 0
	}
}
