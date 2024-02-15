package client

import (
	"log"

	"github.com/msecommerce/cart_service/pkg/config"
	"github.com/sidsreeram/msproto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	Client pb.ProductServiceClient
}

func NewProductClient(c config.Config) *ProductClient {
	cc, err := grpc.Dial(c.ProductSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("client connection failiure", err)
	}
	log.Println("connection success")
	return &ProductClient{
		Client: pb.NewProductServiceClient(cc),
	}
}
