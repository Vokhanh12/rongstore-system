package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	iampb "server/api/iam/v1"
	iam_di "server/internal/iam"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	iamService, err := iam_di.InitializeIamHandler()
	if err != nil {
		log.Fatalf("failed to initialize IAM handler: %v", err)
	}

	iampb.RegisterIamServiceServer(s, iamService)

	log.Println("gRPC server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
