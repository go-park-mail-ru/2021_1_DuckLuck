package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/promo_code"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	proto "github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/proto/cart"
)

type OrderUseCase struct {
	OrderRepo     order.Repository
	CartClient    proto.CartServiceClient
	ProductRepo   product.Repository
	UserRepo      user.Repository
	PromoCodeRepo promo_code.Repository
}

func NewUseCase(orderRepo order.Repository, cartClient proto.CartServiceClient,
	productRepo product.Repository, userRepo user.Repository, promoCodeRepo promo_code.Repository) order.UseCase {
	return &OrderUseCase{
		OrderRepo:     orderRepo,
		CartClient:    cartClient,
		ProductRepo:   productRepo,
		UserRepo:      userRepo,
		PromoCodeRepo: promoCodeRepo,
	}
}

func (u *OrderUseCase) GetPreviewOrder(userId uint64,
	previewCart *models.PreviewCart) (*models.PreviewOrder, error) {
	// Get all info about product in cart
	previewOrder := &models.PreviewOrder{}
	for _, item := range previewCart.Products {
		previewOrder.Products = append(previewOrder.Products,
			&models.PreviewOrderedProducts{
				Id:           item.Id,
				PreviewImage: item.PreviewImage,
			})
	}
	previewOrder.Price = previewCart.Price

	// Get info about user account for order
	userProfile, err := u.UserRepo.SelectProfileById(userId)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	previewOrder.Recipient = models.OrderRecipient{
		FirstName: userProfile.FirstName,
		LastName:  userProfile.LastName,
		Email:     userProfile.Email,
	}

	return previewOrder, nil
}

func (u *OrderUseCase) AddCompletedOrder(order *models.Order, userId uint64,
	previewCart *models.PreviewCart) (*models.OrderNumber, error) {
	if order.PromoCode != "" {
		err := u.PromoCodeRepo.CheckPromo(order.PromoCode)
		if err != nil {
			return nil, errors.ErrPromoCodeNotFound
		}
	}

	price := &models.TotalPrice{}
	for _, product := range previewCart.Products {
		promoPrice, err := u.PromoCodeRepo.GetDiscountPriceByPromo(product.Id, order.PromoCode)
		if err != nil {
			return nil, errors.ErrProductNotFound
		}
		price.TotalBaseCost += promoPrice.BaseCost
		price.TotalCost += promoPrice.TotalCost
	}
	price.TotalDiscount = price.TotalBaseCost - price.TotalCost

	products := previewCart.Products

	orderNumber, err := u.OrderRepo.AddOrder(order, userId, products, price)
	if err != nil {
		return nil, errors.ErrInternalError
	}

	if _, err = u.CartClient.DeleteCart(context.Background(),
		&proto.ReqUserId{UserId: userId}); err != nil {
		return nil, errors.ErrCartNotFound
	}

	return orderNumber, nil
}

func (u *OrderUseCase) GetRangeOrders(userId uint64, paginator *models.PaginatorOrders) (*models.RangeOrders, error) {
	if paginator.PageNum < 1 || paginator.Count < 1 {
		return nil, errors.ErrIncorrectPaginator
	}

	// Max count pages in catalog
	countPages, err := u.OrderRepo.GetCountPages(userId, paginator.Count)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Keys for sort items in catalog
	sortString, err := u.OrderRepo.CreateSortString(paginator.SortKey, paginator.SortDirection)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Get range of products
	orders, err := u.OrderRepo.SelectRangeOrders(userId, sortString, paginator)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Get product for this order
	for _, item := range orders {
		products, err := u.OrderRepo.GetProductsInOrder(item.Id)
		if err != nil {
			return nil, errors.ErrInternalError
		}

		item.Products = products
	}

	return &models.RangeOrders{
		ListPreviewOrders: orders,
		MaxCountPages:     countPages,
	}, nil
}
