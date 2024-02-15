package adapter

import (
	"errors"
	"fmt"
	"log"

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
	queryyy:="SELECT * FROM carts WHERE user_id = ?"
	 err := c.DB.Raw(queryyy,userid).Scan(&cart).Error
	 if err !=nil{
		return models.CartItems{},fmt.Errorf("error in fetching cart id of user :%v",err)
	 } 

	var res models.CartItems
	query := "INSERT INTO cart_items (cart_id, product_id, quantity) VALUES ($1,$2,$3) RETURNING product_id,quantity"

	tx := c.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err = c.DB.Raw(query, cart.ID, req.ProductID, req.Quantity).Scan(&res).Error
	if err != nil {
		tx.Rollback()
		return models.CartItems{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}

	return res, nil

}

func (c *CartAdapter) GetCart(userid uint64) ([]models.CartItems, error) {
	var cart models.Cart
	query1:=`SELECT * FROM carts where user_id=?`
	err:=c.DB.Raw(query1,userid).Scan(&cart).Error
	if err!=nil{
        return []models.CartItems{},fmt.Errorf("error in getting cart details from user_id :%w",err)
	}
	var cartItems []models.CartItems
	tx := c.DB.Begin()
	query:=`SELECT * FROM cart_items where cart_id=?`
	if err := tx.Raw(query,cart.ID).Scan(&cartItems).Error ;err!=nil{
		tx.Rollback()
		return []models.CartItems{},fmt.Errorf("error in fetching cart items :%w",err)
	}
   return cartItems ,nil
}
func (c *CartAdapter) DeleteCartItem(req models.CartItems, userid uint64)(models.CartItems, error) {
	var cart models.Cart
	query1:=`SELECT * FROM carts where user_id=?`
	err:=c.DB.Raw(query1,userid).Scan(&cart).Error
	if err!=nil{
        return models.CartItems{}, fmt.Errorf("error in getting cart details from user_id :%w",err)
	}
	
	tx := c.DB.Begin()
	var qty uint64
	findpro:=`SELECT quantity FROM cart_items where cart_id =$1 and product_id=$2`
	if err := tx.Raw(findpro,cart.ID,req.ProductID).Scan(&qty).Error;err !=nil{
		return models.CartItems{},fmt.Errorf("error in getting quantity of cart")
	}
	log.Println(qty)
	if qty == 0 {
		tx.Rollback()
		return models.CartItems{},fmt.Errorf("no items in cart to reomve")
	}
  var cartItem models.CartItems
	
  if qty == 1 {
	// Delete item from cart if quantity is 1
	dltItem := `DELETE FROM cart_items WHERE cart_id=$1 AND product_id=$2 RETURNING product_id , quantity`
	err := tx.Raw(dltItem, req.CartID).Scan(&cartItem).Error
	if err != nil {
		tx.Rollback()
		return models.CartItems{},  err
	}
} else { 
	// If there is more than one product, reduce the qty by 1
	updateQty := `UPDATE cart_items SET quantity = quantity-1 WHERE cart_id=$1 AND product_id=$2 RETURNING product_id , quantity`
	err = tx.Raw(updateQty, req.CartID, req.ProductID).Scan(&cartItem).Error
	if err != nil {
		tx.Rollback()
		return models.CartItems{}, err
	}
}

	log.Println(cartItem.Quantity)
return models.CartItems{ProductID: cartItem.ProductID,Quantity: cartItem.Quantity},nil

}
func (c *CartAdapter) IncrementQuantity(req models.CartItems, userid uint64) (models.CartItems, error) {
	var cart models.Cart
	tx := c.DB.Begin()

	if err := tx.Where("user_id = ?", userid).First(&cart).Error; err != nil {
		tx.Rollback()
		return req, err
	}

	if req.CartID != uint64(cart.ID) {
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
