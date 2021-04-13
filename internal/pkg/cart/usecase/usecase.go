package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type CartUseCase struct {
	CartRepo    cart.Repository
	ProductRepo product.Repository
}

func NewUseCase(cartRepo cart.Repository, productRepo product.Repository) cart.UseCase {
	return &CartUseCase{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}
}

// Add product in user cart
func (u *CartUseCase) AddProduct(userId uint64, cartArticle *models.CartArticle) error {
	userCart, err := u.CartRepo.SelectCartById(userId)
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

	return u.CartRepo.AddCart(userId, userCart)
}

// Delete product from cart
func (u *CartUseCase) DeleteProduct(userId uint64, identifier *models.ProductIdentifier) error {
	userCart, err := u.CartRepo.SelectCartById(userId)
	if err != nil {
		return err
	}

	// Delete cart of current user
	if len(userCart.Products) == 1 {
		if err = u.CartRepo.DeleteCart(userId); err != nil {
			return err
		}
	}

	delete(userCart.Products, identifier.ProductId)
	return u.CartRepo.AddCart(userId, userCart)
}

// Change product in user cart
func (u *CartUseCase) ChangeProduct(userId uint64, cartArticle *models.CartArticle) error {
	userCart, err := u.CartRepo.SelectCartById(userId)
	if err != nil {
		return err
	}

	if _, ok := userCart.Products[cartArticle.ProductId]; !ok {
		return errors.ErrProductNotFoundInCart
	}
	userCart.Products[cartArticle.ProductId] = &cartArticle.ProductPosition

	return u.CartRepo.AddCart(userId, userCart)
}

// Get preview cart
func (u *CartUseCase) GetPreviewCart(userId uint64) (*models.PreviewCart, error) {
	previewUserCart := &models.PreviewCart{}
	userCart, err := u.CartRepo.SelectCartById(userId)
	switch err {
	case errors.ErrCartNotFound:
		return previewUserCart, err
	case nil:

	default:
		return nil, err
	}

	for id, productPosition := range userCart.Products {
		productById, err := u.ProductRepo.SelectProductById(id)
		if err != nil {
			return nil, err
		}

		previewUserCart.Products = append(previewUserCart.Products,
			&models.PreviewCartArticle{
				Id:           productById.Id,
				Title:        productById.Title,
				Price:        productById.Price,
				PreviewImage: productById.Images[0],
				Count:        productPosition.Count,
			})

		previewUserCart.Price.TotalBaseCost += productById.Price.BaseCost * int(productPosition.Count)
		previewUserCart.Price.TotalDiscount += int(float64(productById.Price.BaseCost) *
			(float64(productById.Price.Discount)) / 100.0 * float64(productPosition.Count))
	}
	previewUserCart.Price.TotalCost = previewUserCart.Price.TotalBaseCost - previewUserCart.Price.TotalDiscount

	return previewUserCart, nil
}

// Delete user cart
func (u *CartUseCase) DeleteCart(userId uint64) error {
	return u.CartRepo.DeleteCart(userId)
}
