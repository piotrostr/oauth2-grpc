package api

import (
	pb "github.com/piotrostr/oauth2-grpc/proto"
)

func fn() {
	_ = pb.Credentials{
		Username: "piotr",
		Password: "secret",
	}
}
