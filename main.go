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

const SwaggerPath string = "./proto/auth.swagger.json"

var ctx = context.Background()

var (
	host              = flag.String("host", "localhost", "host")
	port              = flag.Int("port", 50051, "port to listen on")
	shouldRunAsClient = flag.Bool("client", false, "run as client")
	shouldRunHttp     = flag.Bool("http", false, "enable http server")
)

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
		log.Fatalln(err)
	}
	log.Println(token)

	// Authenticate
	token, err = client.Authenticate(ctx, userDetails.Credentials)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(token)

	// Check if login fails with false credentials
	_, err = client.Authenticate(ctx, &pb.Credentials{
		Username: "piotrostr",
		Password: "wrongpassword",
	})
	if err != nil {
		log.Println(err)
	}
}

func swaggerHandlerFunc(
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

// Create reverse proxy entrypoint, bind grpcServer and serve RPC and HTTP
func runHttp(addr string) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := pb.RegisterAuthServiceHandlerFromEndpoint(
		ctx,
		mux,
		addr,
		opts,
	)
	if err != nil {
		msg := "Error registering AuthServiceHandlerClient: %v"
		log.Fatalf(msg, err)
	}

	addr = ":8081"
	log.Printf("Serving HTTP on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Error while serving: %v", err)
	}
}

func run(addr string) {
	// Grab yourself a port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	authService := api.NewAuthService()
	grpcServer := api.NewGRPCServer()

	pb.RegisterAuthServiceServer(grpcServer, authService)
	reflection.Register(grpcServer)
	log.Println("Serving gRPC")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Error while serving: %v", err)
	}
}

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Grabbing %s for gRPC\n", addr)

	if *shouldRunAsClient {
		runClient(addr)
		return
	}

	if *shouldRunHttp {
		runHttp(addr)
		return
	}

	run(addr)
}
