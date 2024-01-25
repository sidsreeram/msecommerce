package models

import (
	"gorm.io/gorm"
)

// Product represents the database model for products.
type Product struct {
    gorm.Model
    Id          uint64 `gorm:"primaryKey"`
    Name        string
    Quantity    uint64
    Price       uint64
    Description string
    InStock     bool
}
type UpdateProduct struct{
    gorm.Model
    Id         uint64
    Price      uint64
    Quantity   uint64
    Increased    bool
}