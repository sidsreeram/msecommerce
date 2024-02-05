package pkg

import (
	"log"
	"net"

	"github.com/msecommerce/cart_service/pkg/config"
	"github.com/msecommerce/cart_service/pkg/service"
	"github.com/sidsreeram/msproto/pb"
	"google.golang.org/grpc"
)

type ServerHTTP struct {
	engine *grpc.Server
}

func NewServerHTTP(cartService *service.CartServiceServer) *ServerHTTP {
	engine := grpc.NewServer()

	pb.RegisterCartServiceServer(engine, cartService)
	return &ServerHTTP{engine: engine}
}

func (s *ServerHTTP) Start(c config.Config) {
	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalln("failed to listen", err)
	}
	if err = s.engine.Serve(lis); err != nil {
		log.Fatalln("failed to serve", err)
	}
}
