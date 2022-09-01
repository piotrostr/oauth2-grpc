package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/piotrostr/oauth2-grpc/api"
	pb "github.com/piotrostr/oauth2-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var ctx = context.Background()

func grpcHandler(
	grpcServer *grpc.Server,
	httpHandler http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var grpcContentType bool = strings.Contains(
			r.Header.Get("Content-Type"),
			"application/grpc",
		)

		if r.ProtoMajor == 2 && grpcContentType {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}

func main() {
	port := flag.Int("port", 50051, "port to listen on")
	shouldRunHttp := flag.Bool("enable-http", false, "enable http server")
	flag.Parse()

	addr := fmt.Sprintf("localhost:%d", *port)

	authService := api.NewAuthService()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening on %s", addr)

	// TODO add TLS (or ATLS)
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(grpcServer, authService)
	reflection.Register(grpcServer)

	if *shouldRunHttp {
		// Add HTTP router with a route for swagger.json specification
		mux := http.NewServeMux()
		swaggerHandler := func(
			w http.ResponseWriter,
			req *http.Request,
		) {
			f, err := os.ReadFile("./proto/auth.swagger.json")
			if err != nil {
				msg := "Error reading swagger file: %v"
				fmt.Fprintf(w, msg, err)
			}
			io.Copy(w, strings.NewReader(string(f)))
		}
		mux.HandleFunc("/swagger.json", swaggerHandler)

		gatewayMux := runtime.NewServeMux()
		dopts := []grpc.DialOption{
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
		}
		err = pb.RegisterAuthServiceHandlerFromEndpoint(
			ctx,
			gatewayMux,
			addr,
			dopts,
		)
		if err != nil {
			msg := "Error registering AuthServiceHandlerClient: %v"
			log.Fatalf(msg, err)
		}
		mux.Handle("/", gatewayMux)

		httpServer := &http.Server{
			Addr:    addr,
			Handler: grpcHandler(grpcServer, mux),
		}
		log.Println("Serving HTTP and gRPC")
		err = httpServer.Serve(listener)
	} else {
		log.Println("Serving gRPC")
		err = grpcServer.Serve(listener)
	}
	if err != nil {
		log.Fatalf("Error while serving: %v", err)
	}
}
