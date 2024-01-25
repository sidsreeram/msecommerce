package interfaces

import "github.com/msecommerce_product_service/pkg/models"

type ProductAdapter interface {
	Add(req models.Product) (models.Product, error)
	Update(req models.UpdateProduct ) (models.Product, error)
	Delete(id uint64) (string,error)
	Get(id uint64) (models.Product, error)
	GetAll() ([]models.Product, error)
	// GetMultiple(reqs []uint64 ) ([]models.Product, error)
}
