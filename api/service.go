package api

import (
	"context"
	"errors"
	"sync"
	"time"

	pb "github.com/piotrostr/oauth2-grpc/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetToken() *pb.Token {
	return &pb.Token{
		AccessToken: "access_token",
		ExpiresAt:   timestamppb.New(time.Now().Add(time.Hour)),
	}
}

// TODO add database connection
type AuthService struct {
	pb.UnimplementedAuthServiceServer

	mu    sync.Mutex
	users map[string]*pb.UserDetails
}

func NewAuthService() *AuthService {
	return &AuthService{
		users: make(map[string]*pb.UserDetails),
	}
}

func (s *AuthService) CreateAccount(
	ctx context.Context,
	userDetails *pb.UserDetails,
) (*pb.Token, error) {
	s.mu.Lock()
	s.users[userDetails.Credentials.Username] = userDetails
	s.mu.Unlock()
	return GetToken(), nil
}

func (s *AuthService) Authenticate(
	ctx context.Context,
	credentials *pb.Credentials,
) (*pb.Token, error) {
	for _, user := range s.users {
		if user.Credentials.Username == credentials.Username &&
			user.Credentials.Password == credentials.Password {
			return GetToken(), nil
		}
	}
	return nil, errors.New("invalid credentials")
}
