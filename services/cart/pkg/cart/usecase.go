package cart

import "github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/pkg/cart UseCase

type UseCase interface {
	AddProduct(userId uint64, cartArticle *models.CartArticle) error
	DeleteProduct(userId uint64, identifier *models.ProductIdentifier) error
	ChangeProduct(userId uint64, cartArticle *models.CartArticle) error
	GetPreviewCart(userId uint64) (*models.Cart, error)
	DeleteCart(userId uint64) error
}
