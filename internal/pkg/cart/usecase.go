package cart

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type UseCase interface {
	AddProduct(userId uint64, cartArticle *models.CartArticle) error
	DeleteProduct(userId uint64, identifier *models.ProductIdentifier) error
	ChangeProduct(userId uint64, cartArticle *models.CartArticle) error
	GetCart(userId uint64) (*models.Cart, error)
}
