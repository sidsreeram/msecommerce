package di

// import (
// 	"github.com/google/wire"
// 	"github.com/msecommerce/user_service/pkg"
// 	"github.com/msecommerce/user_service/pkg/adapter"
// 	"github.com/msecommerce/user_service/pkg/config"
// 	"github.com/msecommerce/user_service/pkg/db"
// 	"github.com/msecommerce/user_service/pkg/service"
// )

// func InitializeAPI(c config.Config) (*pkg.ServerHTTP, error) {
// 	wire.Build(db.ConnectDatabase,
// 		adapter.NewUserAdapter,
// 		service.NewUserService,
// 		pkg.NewServerHTTP)
// 	return &pkg.ServerHTTP{}, nil
// }
