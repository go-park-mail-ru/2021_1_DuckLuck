package cart

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type UseCase interface {
	AddProduct(userId uint64, position *models.ProductPosition) error
	DeleteProduct(userId uint64, productId uint64) error
	ChangeProduct(userId uint64, position *models.ProductPosition) error
	GetCart(userId uint64) (*models.Cart, error)
}
