package adapter

import (
	"github.com/msecommerce_product_service/pkg/interfaces"
	"github.com/msecommerce_product_service/pkg/models"
	"gorm.io/gorm"
)

type ProductDatabase struct {
	DB *gorm.DB
}

func NewProductAdapter(DB *gorm.DB) interfaces.ProductAdapter {
	return &ProductDatabase{DB}
}
func (p *ProductDatabase) Add(req models.Product)(models.Product,error){
	var pro models.Product
	query := "INSERT INTO products (name,price,quantity,description,in_stock) Values($1,$2,$3,$4,$5) RETURNING name,price,quantity,description,in_stock"
	return pro,p.DB.Raw(query,req.Name,req.Price,req.Quantity,req.Description,req.InStock).Scan(&pro).Error
}
func (p *ProductDatabase) Get(id uint64) (models.Product, error) {
	var pro models.Product
	err := p.DB.First(&pro, id).Error
	if err != nil {
		return models.Product{}, err
	}
	return models.Product{
		Id:          pro.Id,
		Name:        pro.Name,
		Quantity:    pro.Quantity,
		Price:       pro.Price,
		Description: pro.Description,
		InStock:     pro.InStock,
	}, nil
}

func (p *ProductDatabase) Update(req models.UpdateProduct ) (models.Product, error) {
	pro := models.Product{
		Id: req.Id,
		Quantity:    req.Quantity,
		Price:       req.Price,
		InStock:     req.Increased,
	}
	err := p.DB.Save(&pro).Error
	if err != nil {
		return models.Product{}, err
	}
	return models.Product{
		Id:          pro.Id,
		Name:        pro.Name,
		Quantity:    pro.Quantity,
		Price:       pro.Price,
		Description: pro.Description,
		InStock:     pro.InStock,
	}, nil
}


func (p *ProductDatabase) Delete(id uint64) (string,error) {
	var pro models.Product
	err := p.DB.Unscoped().Delete(&pro, id).Error
	if err != nil {
		return "", err
	}
	return "Product Deleted Successfuly",nil
}

func (p *ProductDatabase) GetAll() ([]models.Product, error) {
	var pros []models.Product
	err := p.DB.Find(&pros).Error
	if err != nil {
		return nil, err
	}
	var res []models.Product
	for _, pro := range pros {
		res = append(res, models.Product{
			Id:          pro.Id,
			Name:        pro.Name,
			Quantity:    pro.Quantity,
			Price:       pro.Price,
			Description: pro.Description,
			InStock:     pro.InStock,
		})
	}
	return res, nil
}

// func (p *ProductDatabase) GetMultiple(reqs []uint64 ) ([]models.Product, error) {
// 	var res []models.Product
// 	for _, req := range reqs {
// 		proRes, err := p.Get(req)
// 		if err != nil {
// 			return nil, err
// 		}
// 		res = append(res, proRes)
// 	}
// 	return res, nil
// }
