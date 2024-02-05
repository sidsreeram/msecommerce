package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID uint64 `gorm:"unique"`
}

type CartItems struct {
	gorm.Model
	CartID    uint64 `gorm:"foreignKey:CartID;references:carts(id)"`
	ProductID uint64
	Quantity  uint64
}
