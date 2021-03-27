package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type CartUseCase struct {
	CartRepo cart.Repository
}

func NewUseCase(repo cart.Repository) cart.UseCase {
	return &CartUseCase{
		CartRepo: repo,
	}
}

func (c *CartUseCase) AddProduct(userId uint64, position *models.ProductPosition) error {
	userCart, err := c.CartRepo.GetCart(userId)
	if err != nil {
		userCart = &models.Cart{}
		userCart.Products = make(map[uint64]*models.ProductPosition, 0)
		userCart.Products[position.ProductId] = position
	} else {
		// If product position already exist then increment counter
		if _, ok := userCart.Products[position.ProductId]; ok {
			userCart.Products[position.ProductId].Count++
		} else {
			userCart.Products[position.ProductId] = position
		}
	}

	return c.CartRepo.AddCart(userId, userCart)
}

func (c *CartUseCase) DeleteProduct(userId uint64, productId uint64) error {
	userCart, err := c.CartRepo.GetCart(userId)
	if err != nil {
		return err
	}

	if len(userCart.Products) == 1 {
		if err = c.CartRepo.DeleteCart(userId); err != nil {
			return err
		}
	}

	delete(userCart.Products, productId)
	return c.CartRepo.AddCart(userId, userCart)
}

func (c *CartUseCase) ChangeProduct(userId uint64, position *models.ProductPosition) error {
	userCart, err := c.CartRepo.GetCart(userId)
	if err != nil {
		return err
	}

	if _, ok := userCart.Products[position.ProductId]; !ok {
		return errors.ErrProductNotFoundInCart
	}
	userCart.Products[position.ProductId] = position

	return c.CartRepo.AddCart(userId, userCart)
}

func (c *CartUseCase) GetCart(userId uint64) (*models.Cart, error) {
	return c.CartRepo.GetCart(userId)
}

