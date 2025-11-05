package orsmadapter

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	iampb "server/api/iam/v1"
	"server/pkg/config"
	"server/pkg/logger"
	"server/pkg/metrics"
	"server/pkg/observability"
)

func main() {
	if err := logger.Init(true); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	zlog := logger.L
	defer zlog.Sync()

	cfg := config.Load()
	zlog.Info("config loaded", zap.String("keycloak_url", cfg.KeycloakURL))

	metrics.Register()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		zlog.Info("metrics endpoint started", zap.String("addr", ":9090"))
		if err := http.ListenAndServe(":9090", nil); err != nil {
			zlog.Fatal("metrics server failed", zap.Error(err))
		}
	}()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			observability.GrpcTraceUnaryInterceptor(),
			observability.UnaryServerInterceptor("iam_service", deps.Store, nil, false),
		),
	)

	reflection.Register(s)
	iampb.RegisterIamServiceServer(s, deps.Handler)

	zlog.Info("gRPC server started", zap.String("addr", ":50051"))
	if err := s.Serve(lis); err != nil {
		zlog.Fatal("failed to serve", zap.Error(err))
	}
}
