package cart

import "github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/pkg/models"

type Repository interface {
	SelectCartById(userId uint64) (*models.Cart, error)
	AddCart(userId uint64, userCart *models.Cart) error
	DeleteCart(userId uint64) error
}
