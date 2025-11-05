package titlegladapter

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"server/internal/title-gl-adapter/domain"
	"server/pkg/config"
)

// Đảm bảo interface domain.TitleGl được implement
var _ domain.TitleGl = (*TitleGlClient)(nil)

// TitleGlClient đại diện cho client kết nối tới TileServer GL
type TitleGlClient struct {
	BaseURL string
	Client  *http.Client
	Config  *config.Config
	Health  string
}

// CheckHealth gọi đến endpoint gốc hoặc styles.json để kiểm tra tình trạng server
func (t *TitleGlClient) CheckHealth() error {
	resp, err := t.Client.Get(t.Health)
	if err != nil {
		return fmt.Errorf("TileServer GL health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("TileServer GL health check failed: status %d, body %s", resp.StatusCode, string(body))
	}
	return nil
}

// GetBaseURL trả về BaseURL của TileServer
func (t *TitleGlClient) GetBaseURL() string {
	return t.BaseURL
}

// NewTitleGlClient khởi tạo client TileServer GL
func NewTitleGlClient(cfg *config.Config) domain.TitleGl {
	baseURL := fmt.Sprintf("http://%s:%d", cfg.TitleGlHost, cfg.TitleGlPort)
	healthURL := fmt.Sprintf("%s/", baseURL) // Kiểm tra root URL
	return &TitleGlClient{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 5 * time.Second},
		Config:  cfg,
		Health:  healthURL,
	}
}

// InitTitleGlClient kiểm tra health nhiều lần trước khi trả về client sẵn sàng
func InitTitleGlClient(cfg *config.Config, maxRetries int, interval time.Duration) (domain.TitleGl, error) {
	client := NewTitleGlClient(cfg)
	for i := 0; i < maxRetries; i++ {
		if err := client.CheckHealth(); err == nil {
			return client, nil
		}
		time.Sleep(interval)
	}
	return nil, fmt.Errorf("TileServer GL not ready after %d retries", maxRetries)
}
