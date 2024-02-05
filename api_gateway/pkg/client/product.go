package client

import (
	"fmt"

	"github.com/msecommerce/api_gateway/pkg/config"
	"github.com/sidsreeram/msproto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
type ProductClient struct{
	Client pb.ProductServiceClient
}
func NewProductServiceClient(c config.Config)*ProductClient{
	cc,err:=grpc.Dial(c.ProductSvcUrl,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err !=nil{
		fmt.Errorf("Can't connect with product client :%w",err)
	}	
	
	return &ProductClient{Client: pb.NewProductServiceClient(cc)}
}