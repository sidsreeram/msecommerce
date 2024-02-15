package db

import (
	"github.com/msecommerce/cart_service/pkg/config"
	"github.com/msecommerce/cart_service/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase(c config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(c.DBUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Cart{},
		models.CartItems{})
	return db, err
}
