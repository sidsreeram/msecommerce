package adapter

import (
	"errors"

	"github.com/msecommerce/cart_service/pkg/adapter/interfaces"
	"github.com/msecommerce/cart_service/pkg/models"
	"gorm.io/gorm"
)

type CartAdapter struct {
	DB *gorm.DB
}

func NewCartAdapter(db *gorm.DB) interfaces.CartAdapter {
	return &CartAdapter{db}
}
func (c *CartAdapter) CreateCart(req models.Cart) (models.Cart, error) {
	tx := c.DB.Begin()

	if err := tx.Create(&req).Error; err != nil {
		tx.Rollback()
		return req, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	tx.Commit()
	return req, nil
}
func (c *CartAdapter) AddToCart(req models.CartItems, userid uint64) (models.CartItems, error) {
	var cart models.Cart
	tx := c.DB.Begin()

	if err := tx.Where("user_id = ?", userid).First(&cart).Error; err != nil {
		tx.Rollback()
		return req, err
	}

	req.CartID = uint64(cart.ID) 
	if err := tx.Create(&req).Error; err != nil {
		tx.Rollback()
		return req, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	tx.Commit()
	return req, nil
}

func (c *CartAdapter) GetCart(userid uint64) ([]models.CartItems, error) {
	var cartItems []models.CartItems
	tx := c.DB.Begin()

	if err := tx.Where("user_id = ?", userid).Find(&cartItems).Error; err != nil {
		tx.Rollback()
		return cartItems, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	tx.Commit()
	return cartItems, nil
}
func (c *CartAdapter) DeleteCartItem(req models.CartItems, userid uint64) error {
	var cart models.Cart
	tx := c.DB.Begin()

	if err := tx.Where("user_id = ?", userid).First(&cart).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("cart_id = ?", cart.ID).Delete(&req).Error; err != nil {
		tx.Rollback()
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	tx.Commit()
	return nil
}
func (c *CartAdapter) IncrementQuantity(req models.CartItems, userid uint64) (models.CartItems, error) {
	var cart models.Cart
	tx := c.DB.Begin()

	if err := tx.Where("user_id = ?", userid).First(&cart).Error; err != nil {
		tx.Rollback()
		return req, err
	}

	if req.CartID != uint64(cart.ID)  {
		return req, errors.New("item does not belong to user's cart")
	}

	req.Quantity++
	if err := tx.Save(&req).Error; err != nil {
		tx.Rollback()
		return req, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	tx.Commit()
	return req, nil
}


func (c *CartAdapter) DecrementQuantity(req models.CartItems, userid uint64) (models.CartItems, error) {
	var cart models.Cart
	tx := c.DB.Begin()

	if err := tx.Where("user_id = ?", userid).First(&cart).Error; err != nil {
		tx.Rollback()
		return req, err
	}

	if req.CartID != uint64(cart.ID) {
		return req, errors.New("item does not belong to user's cart")
	}

	if req.Quantity > 0 {
		req.Quantity--
		if err := tx.Save(&req).Error; err != nil {
			tx.Rollback()
			return req, err
		}
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	tx.Commit()
	return req, nil
}
