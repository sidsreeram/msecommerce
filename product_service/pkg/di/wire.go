package di

// import (
// 	"github.com/google/wire"
// 	"github.com/msecommerce_product_service/pkg"
// 	"github.com/msecommerce_product_service/pkg/adapter"
// 	"github.com/msecommerce_product_service/pkg/config"
// 	"github.com/msecommerce_product_service/pkg/db"
// 	"github.com/msecommerce_product_service/pkg/usecase"
// )

// func InitializeAPI(c config.Config) (*pkg.ServerHTTP, error) {
// 	wire.Build(db.ConnectDatabase,
// 		adapter.NewProductAdapter,
// 		usecase.NewProductUsecase,
// 		pkg.NewServerHTTP)

// 	return &pkg.ServerHTTP{}, nil
// }
