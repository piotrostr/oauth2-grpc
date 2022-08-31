package api

import (
	"context"
	"sync"

	pb "github.com/piotrostr/oauth2-grpc/proto"
)

// TODO add database connection
type AuthService struct {
	pb.UnimplementedAuthServiceServer

	mu          sync.Mutex
	userDetails map[string]*pb.UserDetails
}

var ctx = context.Background()

func NewAuthService() *AuthService {
	return &AuthService{
		userDetails: make(map[string]*pb.UserDetails),
	}
}

// TODO implement
func (s *AuthService) CreateAccount(ctx context.Context, userDetails *pb.UserDetails) (*pb.Token, error) {
	return nil, nil
}

// TODO implement
func (s *AuthService) Authenticate(ctx context.Context, credentials *pb.Credentials) (*pb.Token, error) {
	return nil, nil
}
