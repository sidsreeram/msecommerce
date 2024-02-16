package adapter

import (
	"fmt"

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
	queryyy := "SELECT * FROM carts WHERE user_id = ?"
	err := c.DB.Raw(queryyy, userid).Scan(&cart).Error
	if err != nil {
		return models.CartItems{}, fmt.Errorf("error in fetching cart id of user: %v", err)
	}

	var res models.CartItems

	// Check if a cart item with the same cart_id and product_id already exists
	existingCartItem := models.CartItems{}
	err = c.DB.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&existingCartItem).Error
	if err != nil {
		// If not exists, insert a new record
		res = req
		res.CartID = uint64(cart.ID)
		if err := c.DB.Create(&res).Error; err != nil {
			return models.CartItems{}, err
		}
	} else {
		// If exists, update the quantity
		existingCartItem.Quantity += req.Quantity
		if err := c.DB.Save(&existingCartItem).Error; err != nil {
			return models.CartItems{}, err
		}
		res = existingCartItem
	}

	return res, nil
}

func (c *CartAdapter) GetCart(userid uint64) ([]models.CartItems, error) {
	var cart models.Cart
	query1 := `SELECT * FROM carts where user_id=?`
	err := c.DB.Raw(query1, userid).Scan(&cart).Error
	if err != nil {
		return []models.CartItems{}, fmt.Errorf("error in getting cart details from user_id :%w", err)
	}
	var cartItems []models.CartItems
	tx := c.DB.Begin()
	query := `SELECT * FROM cart_items where cart_id=?`
	if err := tx.Raw(query, cart.ID).Scan(&cartItems).Error; err != nil {
		tx.Rollback()
		return []models.CartItems{}, fmt.Errorf("error in fetching cart items :%w", err)
	}
	return cartItems, nil
}
func (c *CartAdapter) DeleteCartItem(req models.CartItems, userid uint64) (models.CartItems, error) {
	var cart models.Cart
	query1 := `SELECT * FROM carts WHERE user_id=?`
	err := c.DB.Raw(query1, userid).Scan(&cart).Error
	if err != nil {
		return models.CartItems{}, fmt.Errorf("error in getting cart details from user_id: %w", err)
	}

	tx := c.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var qty uint64
	findpro := `SELECT quantity FROM cart_items WHERE cart_id = $1 AND product_id = $2`
	if err := tx.Raw(findpro, cart.ID, req.ProductID).Scan(&qty).Error; err != nil {
		tx.Rollback()
		return models.CartItems{}, fmt.Errorf("error in getting quantity of cart: %w", err)
	}

	if qty == 0 {
		tx.Rollback()
		return models.CartItems{}, fmt.Errorf("no items in cart to remove")
	}

	var cartItem models.CartItems

	if qty == 1 {
		// Delete item from cart if quantity is 1
		dltItem := `DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2 RETURNING product_id, quantity`
		err := tx.Raw(dltItem, cart.ID, req.ProductID).Scan(&cartItem).Error
		if err != nil {
			tx.Rollback()
			return models.CartItems{}, err
		}
	} else {
		// If there is more than one product, reduce the qty by 1
		updateQty := `UPDATE cart_items SET quantity = quantity - 1 WHERE cart_id = $1 AND product_id = $2 RETURNING product_id, quantity`
		err = tx.Raw(updateQty, cart.ID, req.ProductID).Scan(&cartItem).Error
		if err != nil {
			tx.Rollback()
			return models.CartItems{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return models.CartItems{}, err
	}

	return models.CartItems{ProductID: cartItem.ProductID, Quantity: cartItem.Quantity}, nil
}

func (c *CartAdapter) IncrementQuantity(req models.CartItems, userid uint64) (models.CartItems, error) {
	var cart models.Cart
	tx := c.DB.Begin()

	if err := tx.Where("user_id = ?", userid).First(&cart).Error; err != nil {
		tx.Rollback()
		return req, err
	}

	var cartItem models.CartItems
	if err := tx.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error; err != nil {
		tx.Rollback()
		return req, err
	}

	cartItem.Quantity++
	if err := tx.Save(&cartItem).Error; err != nil {
		tx.Rollback()
		return req, err
	}

	tx.Commit()
	return cartItem, nil
}

func (c *CartAdapter) DecrementQuantity(req models.CartItems, userid uint64) (models.CartItems, error) {
	var cart models.Cart
	tx := c.DB.Begin()

	if err := tx.Where("user_id = ?", userid).First(&cart).Error; err != nil {
		tx.Rollback()
		return req, err
	}

	var cartItem models.CartItems
	if err := tx.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error; err != nil {
		tx.Rollback()
		return req, err
	}

	if cartItem.Quantity > 0 {
		cartItem.Quantity--
		if err := tx.Save(&cartItem).Error; err != nil {
			tx.Rollback()
			return req, err
		}
	}

	tx.Commit()
	return cartItem, nil
}
