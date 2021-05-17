package promo_code

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type Repository interface {
	GetDiscountPriceByPromo(productId uint64, promoCode string) (*models.ProductPrice, error)
}