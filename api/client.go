package api

import (
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
	defer conn.Close()
	return &Client{
		conn:       conn,
		authClient: pb.NewAuthServiceClient(conn),
	}
}

func (c *Client) Authenticate(credentials *pb.Credentials) {
	c.Authenticate(credentials)
}

func (c *Client) CreateAccount(userDetails *pb.UserDetails) {
	c.CreateAccount(userDetails)
}
