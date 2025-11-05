package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	titlegl "server/internal/title-gl-adapter/infrastructure/client"
	"server/pkg/config"
	"server/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// --- Init logger ---
	if err := logger.Init(true); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	zlog := logger.L
	defer zlog.Sync()

	// --- Load config ---
	cfg := config.Load()
	zlog.Info("✅ Loaded config",
		zap.String("title_gl_host", cfg.TitleGlHost),
		zap.Int("title_gl_port", cfg.TitleGlPort),
	)

	// --- Initialize Title GL client ---
	client, err := titlegl.InitTitleGlClient(cfg, 10, 3*time.Second)
	if err != nil {
		zlog.Fatal("❌ Failed to connect to Title GL", zap.Error(err))
	}

	zlog.Info("✅ Title GL client ready",
		zap.String("base_url", client.GetBaseURL()),
	)

	// --- Expose health endpoint ---
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := client.CheckHealth(); err != nil {
			http.Error(w, fmt.Sprintf("TitleGL unhealthy: %v", err), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TitleGL OK"))
	})

	addr := ":8085" // có thể đổi port
	zlog.Info("🚀 Title GL adapter running", zap.String("addr", addr))
	if err := http.ListenAndServe(addr, nil); err != nil {
		zlog.Fatal("server failed", zap.Error(err))
	}
}
