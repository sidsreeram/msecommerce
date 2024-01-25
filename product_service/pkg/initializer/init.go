package initializer

import (
	"github.com/msecommerce_product_service/pkg/adapter"
	"github.com/msecommerce_product_service/pkg/usecase"
	"gorm.io/gorm"
)
func Initialize(db *gorm.DB) *usecase.ProductUsecase {

	adapter := adapter.NewProductAdapter(db)
	usecase := usecase.NewProductUsecase(adapter)

	return usecase
}
