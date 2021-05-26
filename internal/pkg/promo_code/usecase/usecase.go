package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/promo_code"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type PromoCodeUseCase struct {
	PromoCodeRepo promo_code.Repository
}

func NewUseCase(promoCodeRepo promo_code.Repository) promo_code.UseCase {
	return &PromoCodeUseCase{
		PromoCodeRepo: promoCodeRepo,
	}
}

func (u *PromoCodeUseCase) ApplyPromoCodeToOrder(promoCodeGroup *models.PromoCodeGroup) (*models.DiscountedPrice, error) {
	err := u.PromoCodeRepo.CheckPromo(promoCodeGroup.PromoCode)
	if err != nil {
		return nil, errors.ErrPromoCodeNotFound
	}

	productsInAction := 0
	discountedPrice := &models.DiscountedPrice{}
	for _, productId := range promoCodeGroup.Products {
		price, err := u.PromoCodeRepo.GetDiscountPriceByPromo(productId, promoCodeGroup.PromoCode)

		if err == nil {
			productsInAction += 1
		} else if err != errors.ErrProductNotInPromo {
			return nil, errors.ErrProductNotFound
		}

		discountedPrice.TotalBaseCost += price.BaseCost
		discountedPrice.TotalCost += price.TotalCost
	}
	if productsInAction == 0 {
		return nil, errors.ErrProductNotInPromo
	}
	discountedPrice.TotalDiscount = discountedPrice.TotalBaseCost - discountedPrice.TotalCost

	return discountedPrice, nil
}
