package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
)

type OrderUseCase struct {
	OrderRepo   order.Repository
	CartRepo    cart.Repository
	ProductRepo product.Repository
	UserRepo    user.Repository
}

func NewUseCase(orderRepo order.Repository, cartRepo cart.Repository,
	productRepo product.Repository, userRepo user.Repository) order.UseCase {
	return &OrderUseCase{
		OrderRepo:   orderRepo,
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
		UserRepo:    userRepo,
	}
}

func (u *OrderUseCase) GetPreviewOrder(userId uint64,
	previewCart *models.PreviewCart) (*models.PreviewOrder, error) {
	// Get all info about product in cart
	previewOrder := &models.PreviewOrder{}
	previewOrder.Products = previewCart.Products
	previewOrder.Price = previewCart.Price

	// Get info about user account for order
	userProfile, err := u.UserRepo.SelectProfileById(userId)
	if err != nil {
		return nil, err
	}
	previewOrder.Recipient = models.OrderRecipient{
		FirstName: userProfile.FirstName.String,
		LastName:  userProfile.LastName.String,
		Email:     userProfile.Email,
	}

	return previewOrder, nil
}

func (u *OrderUseCase) AddCompletedOrder(order *models.Order, userId uint64,
	previewCart *models.PreviewCart) (uint64, error) {
	// Get all info about product in cart
	products := previewCart.Products
	price := previewCart.Price

	orderId, err := u.OrderRepo.AddOrder(order, userId, products, &price)
	if err != nil {
		return 0, err
	}

	if err = u.CartRepo.DeleteCart(userId); err != nil {
		return 0, err
	}

	return orderId, nil
}
