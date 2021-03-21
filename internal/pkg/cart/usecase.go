package cart

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type UseCase interface {
	AddProductInCart(userId uint64, position *models.ProductPosition) error
	DeleteProductInCart(userId uint64, productId uint64) error
	ChangeProductInCart(userId uint64, position *models.ProductPosition) error
	GetProductsFromCart(userId uint64) (*models.Cart, error)
}
