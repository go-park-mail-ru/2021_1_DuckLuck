package cart

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type Repository interface {
	AddProductPosition(userId uint64, position *models.ProductPosition) error
	DeleteProductPosition(userId uint64, productId uint64) error
	UpdateProductPosition(userId uint64, position *models.ProductPosition) error
	GetProductsFromCart(userId uint64) (*models.Cart, error)
	AddProductsInCart(userId uint64, userCart *models.Cart) error
}
