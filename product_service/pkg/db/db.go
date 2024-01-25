package db

import (
	"github.com/msecommerce_product_service/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(Connect string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(Connect), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Product{})
	return db, err
}
