package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/piotrostr/oauth2-grpc/api"
	pb "github.com/piotrostr/oauth2-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	Host        string = "localhost"
	SwaggerPath string = "./proto/auth.swagger.json"
)

var ctx = context.Background()

func grpcHandler(
	grpcServer *grpc.Server,
	httpHandler http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(
			r.Header.Get("Content-Type"),
			"application/grpc",
		) {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}

// runClient runs a test client to check if grpc server works as intended
func runClient(addr string) {
	client := api.NewClient(addr)
	userDetails := &pb.UserDetails{
		Credentials: &pb.Credentials{
			Username: "piotrostr",
			Password: "password",
		},
	}

	// Create account (overwrite if exists)
	token, err := client.CreateAccount(ctx, userDetails)
	if err != nil {
		glog.Fatalln(err)
	}
	glog.Infoln(token)

	// Authenticate
	token, err = client.Authenticate(ctx, userDetails.Credentials)
	if err != nil {
		glog.Fatalln(err)
	}
	glog.Infoln(token)

	// Check if glogin fails with false credentials
	_, err = client.Authenticate(ctx, &pb.Credentials{
		Username: "piotrostr",
		Password: "wrongpassword",
	})
	if err != nil {
		glog.Infoln(err)
	}
}

func runHttp(addr string, grpcServer *grpc.Server, listener net.Listener) {
	// Add HTTP router with a route for swagger.json specification
	mux := http.NewServeMux()
	swaggerHandler := func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		f, err := os.ReadFile(SwaggerPath)
		if err != nil {
			msg := "Error reading swagger file: %v"
			fmt.Fprintf(w, msg, err)
		}
		_, err = io.Copy(w, strings.NewReader(string(f)))
		if err != nil {
			msg := "Error copying over swagger file: %v"
			fmt.Fprintf(w, msg, err)
		}
	}
	mux.HandleFunc("/swagger.json", swaggerHandler)

	gatewayMux := runtime.NewServeMux()
	dopts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := pb.RegisterAuthServiceHandlerFromEndpoint(
		ctx,
		gatewayMux,
		addr,
		dopts,
	)
	if err != nil {
		msg := "Error registering AuthServiceHandlerClient: %v"
		glog.Fatalf(msg, err)
	}
	mux.Handle("/", gatewayMux)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: grpcHandler(grpcServer, mux),
	}
	glog.Infoln("Serving HTTP and gRPC")
	err = httpServer.Serve(listener)
	if err != nil {
		glog.Fatalf("Error while serving: %v", err)
	}
}

func main() {
	port := flag.Int("port", 50051, "port to listen on")
	shouldRunAsClient := flag.Bool("client", false, "run as client")
	shouldRunHttp := flag.Bool("http", false, "enable http server")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", Host, *port)

	if *shouldRunAsClient {
		runClient(addr)
		return
	}

	// grab yourself a port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatalf("Failed to listen: %v", err)
	}
	glog.Infof("Listening on %s", addr)

	authService := api.NewAuthService()
	grpcServer := api.NewGRPCServer()

	pb.RegisterAuthServiceServer(grpcServer, authService)
	reflection.Register(grpcServer)

	// TODO this doesn't work due to a certificate error (including the
	// example implementation on github.com/philips/grpc-gateway)
	if *shouldRunHttp {
		runHttp(addr, grpcServer, listener)
	}

	glog.Infoln("Serving gRPC")
	err = grpcServer.Serve(listener)
	if err != nil {
		glog.Fatalf("Error while serving: %v", err)
	}
}
