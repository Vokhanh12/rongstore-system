package main

import (
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	iampb "server/api/iam/v1"
	wire "server/internal/iam"
	"server/pkg/errors"
	"server/pkg/logger"
	"server/pkg/metrics"
	"server/pkg/observability"
)

func main() {
	// 0) Init logger
	if err := logger.Init(true); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	zlog := logger.L
	defer zlog.Sync()

	// 1) Load error catalog
	catalogPath := filepath.Join(".", "errors.yaml")
	if err := errors.InitFromYAML(catalogPath); err != nil {
		zlog.Warn("cannot load errors.yaml, using static mapping", zap.String("path", catalogPath), zap.Error(err))
	}

	// 2) Prometheus metrics
	metrics.Register()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		zlog.Info("metrics endpoint started", zap.String("addr", ":9090"))
		if err := http.ListenAndServe(":9090", nil); err != nil {
			zlog.Fatal("metrics server failed", zap.Error(err))
		}
	}()

	// 3) Listen gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		zlog.Fatal("failed to listen", zap.Error(err))
	}

	// 4) Wire up
	deps, err := wire.InitializeIamHandler()
	if err != nil {
		zlog.Fatal("failed to initialize IAM deps", zap.Error(err))
	}

	// 5) gRPC server + interceptors (trace + obs with session store)
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			observability.GrpcTraceUnaryInterceptor(),
			observability.UnaryServerInterceptor("iam_service", deps.Store),
		),
	)

	reflection.Register(s)
	iampb.RegisterIamServiceServer(s, deps.Handler)

	zlog.Info("gRPC server started", zap.String("addr", ":50051"))
	if err := s.Serve(lis); err != nil {
		zlog.Fatal("failed to serve", zap.Error(err))
	}
}
