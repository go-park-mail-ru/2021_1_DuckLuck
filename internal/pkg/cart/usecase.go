package cart

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type UseCase interface {
	AddProductInCart(userId uint64, position *models.ProductPosition) error
}
