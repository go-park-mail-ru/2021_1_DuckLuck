package promo_code

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type UseCase interface {
	ApplyPromoCodeToOrder(promoCodeGroup *models.PromoCodeGroup) (*models.DiscountedPrice, error)
}
