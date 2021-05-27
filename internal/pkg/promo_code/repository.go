package promo_code

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/promo_code Repository

type Repository interface {
	GetDiscountPriceByPromo(productId uint64, promoCode string) (*models.PromoPrice, error)
	CheckPromo(promoCode string) error
}
