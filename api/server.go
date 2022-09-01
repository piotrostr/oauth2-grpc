package api

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewGRPCServer() *grpc.Server {
	// TODO have a look at ATLS as well as an alternative to TLS below
	creds, err := credentials.NewServerTLSFromFile(
		CertificatePath,
		KeyPath,
	)
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	return grpcServer
}
