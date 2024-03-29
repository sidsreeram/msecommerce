package db

import (
	"log"

	"github.com/msecommerce_product_service/pkg/config"
	"github.com/msecommerce_product_service/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(c config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(c.DBUrl), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.Product{})
	return db, err
}
