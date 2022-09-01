package api

import (
	"context"
	"log"

	pb "github.com/piotrostr/oauth2-grpc/proto"
	"google.golang.org/grpc"
)

type Client struct {
	conn       *grpc.ClientConn
	authClient pb.AuthServiceClient
}

func NewClient(url string) *Client {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	return &Client{
		conn:       conn,
		authClient: pb.NewAuthServiceClient(conn),
	}
}

func (c *Client) Authenticate(
	ctx context.Context,
	credentials *pb.Credentials,
) (*pb.Token, error) {
	return c.authClient.Authenticate(ctx, credentials)
}

func (c *Client) CreateAccount(
	ctx context.Context,
	userDetails *pb.UserDetails,
) (*pb.Token, error) {
	return c.authClient.CreateAccount(ctx, userDetails)
}
