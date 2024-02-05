package pkg

import (
	"log"
	"net"
	"github.com/sidsreeram/msproto/pb"

	"github.com/msecommerce/user_service/pkg/config"
	"github.com/msecommerce/user_service/pkg/service"
	"google.golang.org/grpc"
)

type ServerHTTP struct {
	engine *grpc.Server
}

func NewServerHTTP(userService *service.UserService) *ServerHTTP {
	engine := grpc.NewServer()

	pb.RegisterUserServiceServer(engine,userService )
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
