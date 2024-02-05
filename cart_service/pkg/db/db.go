package db

import (
	"github.com/msecommerce/cart_service/pkg/config"
	"github.com/msecommerce/cart_service/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(c config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(c.DBUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Cart{},
		models.CartItems{})
	return db, err
}
