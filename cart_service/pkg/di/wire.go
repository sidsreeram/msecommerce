package di

// import (
// 	"github.com/google/wire"
// 	"github.com/msecommerce/cart_service/pkg"
// 	"github.com/msecommerce/cart_service/pkg/adapter"
// 	"github.com/msecommerce/cart_service/pkg/client"
// 	"github.com/msecommerce/cart_service/pkg/config"
// 	"github.com/msecommerce/cart_service/pkg/db"
// 	"github.com/msecommerce/cart_service/pkg/service"
// )

// func InitializeAPI(c config.Config) (*pkg.ServerHTTP, error) {
// 	wire.Build(db.ConnectDatabase,
// 		adapter.NewCartAdapter,
// 	    client.NewProductClient,
// 		service.NewCartServiceServer,
// 		pkg.NewServerHTTP)

// 	return &pkg.ServerHTTP{}, nil
// }
