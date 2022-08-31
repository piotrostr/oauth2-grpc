package main

import (
	"log"
	"net"

	"github.com/piotrostr/oauth2-grpc/api"
	pb "github.com/piotrostr/oauth2-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	authService := api.NewAuthService()
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(grpcServer, authService)
	reflection.Register(grpcServer)
	log.Println("Serving on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error while serving: %v", err)
	}
	// TODO use ServeHTTP instead of Serve to provide both REST and gRPC
}
