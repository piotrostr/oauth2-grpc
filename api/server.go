package api

import (
	"crypto/tls"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	CertificatePath string = "./certs/server.crt"
	KeyPath         string = "./certs/server.key"
	UseTLS          bool   = false
)

// TODO have a look at ATLS as well as an alternative to TLS below
func NewGRPCServer() *grpc.Server {
	if UseTLS {
		cert, err := tls.LoadX509KeyPair(CertificatePath, KeyPath)
		if err != nil {
			log.Fatalf("could not load tls cert: %v", err)
		}
		creds := credentials.NewServerTLSFromCert(&cert)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		return grpc.NewServer(grpc.Creds(creds))
	} else {
		return grpc.NewServer()
	}
}
