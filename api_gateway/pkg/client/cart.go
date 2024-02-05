package client

import (
	"fmt"

	"github.com/msecommerce/api_gateway/pkg/config"
	"github.com/sidsreeram/msproto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CartClient struct {
	Client pb.CartServiceClient
}

func NewCartClient(c config.Config) *CartClient {
	cc, err := grpc.Dial(c.CartSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("can't connect with cart client :%w", err)
	}
	fmt.Println("hiii")
	return &CartClient{
		Client: pb.NewCartServiceClient(cc),
	}
}
