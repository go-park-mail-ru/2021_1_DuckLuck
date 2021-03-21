package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
)

type CartUseCase struct {
	CartRepo cart.Repository
}

func NewUseCase(repo cart.Repository) cart.UseCase {
	return &CartUseCase{
		CartRepo: repo,
	}
}

func (c *CartUseCase) AddProductInCart(userId uint64, position *models.ProductPosition) error {
	return c.CartRepo.AddProductPosition(userId, position)
}
