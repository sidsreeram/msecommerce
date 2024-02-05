package client

import (
	"fmt"

	"github.com/msecommerce/api_gateway/pkg/config"
	"github.com/sidsreeram/msproto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
type UserClient struct {
	Client pb.UserServiceClient
}
func NewUserClient(c config.Config) (*UserClient, error) {
	cc, err := grpc.Dial(c.UserSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Could not connect with user client :%w", err)
	}
	return &UserClient{
		Client: pb.NewUserServiceClient(cc),
	}, nil
}
