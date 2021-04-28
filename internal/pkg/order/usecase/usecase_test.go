package usecase

import (
	"testing"

	cart_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	order_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order/mock"
	product_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product/mock"
	user_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserUseCase_GetPreviewOrder(t *testing.T) {
	userId := uint64(3)
	profileUser := models.ProfileUser{
		Id:        1,
		FirstName: "test",
		LastName:  "last",
		Email:     "test@test.ru",
		Password:  nil,
		Avatar: models.Avatar{
			Url: "httt://test.png",
		},
	}
	previewOrder := models.PreviewOrder{
		Products: nil,
		Recipient: models.OrderRecipient{
			FirstName: "test",
			LastName:  "last",
			Email:     "test@test.ru",
		},
		Price: models.TotalPrice{
			TotalDiscount: 32,
			TotalCost:     234,
			TotalBaseCost: 34,
		},
		Address: models.OrderAddress{
			Address: "",
		},
	}
	previewCart := models.PreviewCart{
		Products: nil,
		Price: models.TotalPrice{
			TotalDiscount: 32,
			TotalCost:     234,
			TotalBaseCost: 34,
		},
	}

	t.Run("GetPreviewOrder_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)
		userRepo.
			EXPECT().
			SelectProfileById(userId).
			Return(&profileUser, nil)
		productRepo := product_mock.NewMockRepository(ctrl)
		orderRepo := order_mock.NewMockRepository(ctrl)
		cartRepo := cart_mock.NewMockRepository(ctrl)

		userUCase := NewUseCase(orderRepo, cartRepo, productRepo, userRepo)

		userData, err := userUCase.GetPreviewOrder(userId, &previewCart)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, previewOrder, *userData, "not equal data")
	})

	t.Run("GetPreviewOrder_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)
		userRepo.
			EXPECT().
			SelectProfileById(userId).
			Return(nil, errors.ErrDBInternalError)
		productRepo := product_mock.NewMockRepository(ctrl)
		orderRepo := order_mock.NewMockRepository(ctrl)
		cartRepo := cart_mock.NewMockRepository(ctrl)

		userUCase := NewUseCase(orderRepo, cartRepo, productRepo, userRepo)

		_, err := userUCase.GetPreviewOrder(userId, &previewCart)
		assert.Equal(t, err, errors.ErrUserNotFound, "not equal errors")
	})
}

func TestUserUseCase_AddCompletedOrder(t *testing.T) {
	userId := uint64(3)
	order := models.Order{
		Recipient: models.OrderRecipient{
			FirstName: "test",
			LastName:  "last",
			Email:     "test@test.ru",
		},
		Address: models.OrderAddress{
			Address: "test",
		},
	}
	previewCart := models.PreviewCart{
		Products: nil,
		Price: models.TotalPrice{
			TotalDiscount: 32,
			TotalCost:     234,
			TotalBaseCost: 34,
		},
	}
	orderId := uint64(4)

	t.Run("GetPreviewOrder_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)
		productRepo := product_mock.NewMockRepository(ctrl)
		orderRepo := order_mock.NewMockRepository(ctrl)
		orderRepo.
			EXPECT().
			AddOrder(&order, userId, previewCart.Products, &previewCart.Price).
			Return(orderId, nil)

		cartRepo := cart_mock.NewMockRepository(ctrl)
		cartRepo.
			EXPECT().
			DeleteCart(userId).
			Return(nil)

		userUCase := NewUseCase(orderRepo, cartRepo, productRepo, userRepo)

		userData, err := userUCase.AddCompletedOrder(&order, userId, &previewCart)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, orderId, userData, "not equal data")
	})

	t.Run("GetPreviewOrder_internal_sever_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)
		productRepo := product_mock.NewMockRepository(ctrl)
		orderRepo := order_mock.NewMockRepository(ctrl)
		orderRepo.
			EXPECT().
			AddOrder(&order, userId, previewCart.Products, &previewCart.Price).
			Return(orderId, errors.ErrDBInternalError)

		cartRepo := cart_mock.NewMockRepository(ctrl)

		userUCase := NewUseCase(orderRepo, cartRepo, productRepo, userRepo)

		_, err := userUCase.AddCompletedOrder(&order, userId, &previewCart)
		assert.Equal(t, err, errors.ErrInternalError, "not equal errors")
	})

	t.Run("GetPreviewOrder_cart_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)
		productRepo := product_mock.NewMockRepository(ctrl)
		orderRepo := order_mock.NewMockRepository(ctrl)
		orderRepo.
			EXPECT().
			AddOrder(&order, userId, previewCart.Products, &previewCart.Price).
			Return(orderId, nil)

		cartRepo := cart_mock.NewMockRepository(ctrl)
		cartRepo.
			EXPECT().
			DeleteCart(userId).
			Return(errors.ErrCartNotFound)

		userUCase := NewUseCase(orderRepo, cartRepo, productRepo, userRepo)

		_, err := userUCase.AddCompletedOrder(&order, userId, &previewCart)
		assert.Equal(t, err, errors.ErrCartNotFound, "not equal errors")
	})
}
