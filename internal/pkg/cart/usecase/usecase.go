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

func (c *CartUseCase) AddProduct(userId uint64, cartArticle *models.CartArticle) error {
	userCart, err := c.CartRepo.GetCart(userId)
	if err != nil {
		userCart = &models.Cart{}
		userCart.Products = make(map[uint64]*models.ProductPosition, 0)
		userCart.Products[cartArticle.ProductId] = &cartArticle.ProductPosition
	} else {
		// If product position already exist then increment counter
		if _, ok := userCart.Products[cartArticle.ProductId]; ok {
			userCart.Products[cartArticle.ProductId].Count += cartArticle.Count
		} else {
			userCart.Products[cartArticle.ProductId] = &cartArticle.ProductPosition
		}
	}

	return c.CartRepo.AddCart(userId, userCart)
}

func (c *CartUseCase) DeleteProduct(userId uint64, identifier *models.ProductIdentifier) error {
	userCart, err := c.CartRepo.GetCart(userId)
	if err != nil {
		return err
	}

	// Delete cart of current user
	if len(userCart.Products) == 1 {
		if err = c.CartRepo.DeleteCart(userId); err != nil {
			return err
		}
	}

	delete(userCart.Products, identifier.ProductId)
	return c.CartRepo.AddCart(userId, userCart)
}

func (c *CartUseCase) ChangeProduct(userId uint64, cartArticle *models.CartArticle) error {
	userCart, err := c.CartRepo.GetCart(userId)
	if err != nil {
		return err
	}

	if _, ok := userCart.Products[cartArticle.ProductId]; !ok {
		return errors.ErrProductNotFoundInCart
	}
	userCart.Products[cartArticle.ProductId] = &cartArticle.ProductPosition

	return c.CartRepo.AddCart(userId, userCart)
}

func (c *CartUseCase) GetCart(userId uint64) (*models.Cart, error) {
	return c.CartRepo.GetCart(userId)
}
