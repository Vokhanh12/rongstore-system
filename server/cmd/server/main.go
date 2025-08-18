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
	iam_di "server/internal/iam"

	"server/pkg/errors"
	"server/pkg/logger"
	"server/pkg/metrics"
	"server/pkg/observability"
)

func main() {
	// 0) init logger (production mode)
	if err := logger.Init(true); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	zlog := logger.L
	defer zlog.Sync()

	// 1) load error catalog (optional)
	catalogPath := filepath.Join(".", "errors.yaml")
	if err := errors.InitFromYAML(catalogPath); err != nil {
		zlog.Warn("cannot load errors.yaml, using static mapping", zap.String("path", catalogPath), zap.Error(err))
	}

	// 2) register prometheus metrics and expose /metrics
	metrics.Register()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		zlog.Info("metrics endpoint started", zap.String("addr", ":9090"))
		if err := http.ListenAndServe(":9090", nil); err != nil {
			zlog.Fatal("metrics server failed", zap.Error(err))
		}
	}()

	// 3) listen grpc
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		zlog.Fatal("failed to listen", zap.Error(err))
	}

	// 4) create gRPC server with interceptors:
	// - trace extractor (metadata -> ctx)
	// - observability (metrics + logging)
	// - translate errors to grpc status
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			observability.GrpcTraceUnaryInterceptor(), // inject trace_id
			//observability.UnaryLoggingInterceptor(),             // logs request
			observability.UnaryServerInterceptor("iam_service"), // metrics & more
		),
	)

	reflection.Register(s)

	// 5) initialize services / handlers
	iamService, err := iam_di.InitializeIamHandler()
	if err != nil {
		zlog.Fatal("failed to initialize IAM handler", zap.Error(err))
	}

	iampb.RegisterIamServiceServer(s, iamService)

	zlog.Info("gRPC server started", zap.String("addr", ":50051"))
	if err := s.Serve(lis); err != nil {
		zlog.Fatal("failed to serve", zap.Error(err))
	}
}
