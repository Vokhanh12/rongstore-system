package auth

type HeaderRule struct {
	PathPrefix string   // ví dụ "/v1/iam/handshake"
	Required   []string // các header bắt buộc
	Optional   []string // có thể thêm nếu muốn
}

var DefaultRules = []HeaderRule{
	{
		PathPrefix: "/v1/iam/handshake",
		Required:   []string{"session_id"},
	},
	{
		PathPrefix: "/",
		Required:   []string{"Authorization"},
	},
}
