package interfaces

import "github.com/msecommerce/cart_service/pkg/models"

type CartAdapter interface {
	CreateCart(req models.Cart) (models.Cart, error)
	AddToCart(req models.CartItems, userid uint64) (models.CartItems, error)
	GetCart(userid uint64) ([]models.CartItems, error)
	DeleteCartItem(req models.CartItems, userid uint64) error
	IncrementQuantity(req models.CartItems, userid uint64) (models.CartItems, error)
	DecrementQuantity(req models.CartItems, userid uint64) (models.CartItems, error)
}
