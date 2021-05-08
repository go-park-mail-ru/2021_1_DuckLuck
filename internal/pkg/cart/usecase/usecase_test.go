package usecase

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	product_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	cart_service "github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/proto/cart"
	proto "github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/proto/cart"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
)

func TestCartUseCase_AddProduct(t *testing.T) {
	userId := uint64(2)
	count := uint64(4)
	productId := uint64(1)
	cartArticle := &models.CartArticle{
		ProductPosition:   models.ProductPosition{
			Count: count,
		},
		ProductIdentifier: models.ProductIdentifier{
			ProductId: productId,
		},
	}
	reqCartArticle := &proto.ReqCartArticle{
		Position:  &proto.ProductPosition{Count: count},
		ProductId: productId,
		UserId:    userId,
	}

	t.Run("AddProduct_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)
		cartClient := cart_service.NewMockCartServiceClient(ctrl)

		cartClient.
			EXPECT().
			AddProduct(context.Background(), reqCartArticle).
			Return(&empty.Empty{}, nil)

		userUCase := NewUseCase(cartClient, productRepo)

		err := userUCase.AddProduct(userId, cartArticle)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("AddProduct_internal_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)
		cartClient := cart_service.NewMockCartServiceClient(ctrl)

		cartClient.
			EXPECT().
			AddProduct(context.Background(), reqCartArticle).
			Return(&empty.Empty{}, errors.ErrInternalError)

		userUCase := NewUseCase(cartClient, productRepo)

		err := userUCase.AddProduct(userId, cartArticle)
		assert.Error(t, err, "expected internal error")
	})
}

func TestCartUseCase_DeleteProduct(t *testing.T) {
	userId := uint64(2)
	productId := uint64(1)
	productIdentifier := &models.ProductIdentifier{
		ProductId: productId,
	}
	reqProductIdentifier := &proto.ReqProductIdentifier{
		ProductId: productId,
		UserId:    userId,
	}

	t.Run("DeleteProduct_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)
		cartClient := cart_service.NewMockCartServiceClient(ctrl)

		cartClient.
			EXPECT().
			DeleteProduct(context.Background(), reqProductIdentifier).
			Return(&empty.Empty{}, nil)

		userUCase := NewUseCase(cartClient, productRepo)

		err := userUCase.DeleteProduct(userId, productIdentifier)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("DeleteProduct_internal_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)
		cartClient := cart_service.NewMockCartServiceClient(ctrl)

		cartClient.
			EXPECT().
			DeleteProduct(context.Background(), reqProductIdentifier).
			Return(&empty.Empty{}, errors.ErrInternalError)

		userUCase := NewUseCase(cartClient, productRepo)

		err := userUCase.DeleteProduct(userId, productIdentifier)
		assert.Error(t, err, "expected internal error")
	})
}

func TestCartUseCase_ChangeProduct(t *testing.T) {
	userId := uint64(2)
	count := uint64(4)
	productId := uint64(1)
	cartArticle := &models.CartArticle{
		ProductPosition:   models.ProductPosition{
			Count: count,
		},
		ProductIdentifier: models.ProductIdentifier{
			ProductId: productId,
		},
	}
	reqCartArticle := &proto.ReqCartArticle{
		Position:  &proto.ProductPosition{Count: count},
		ProductId: productId,
		UserId:    userId,
	}

	t.Run("ChangeProduct_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)
		cartClient := cart_service.NewMockCartServiceClient(ctrl)

		cartClient.
			EXPECT().
			ChangeProduct(context.Background(), reqCartArticle).
			Return(&empty.Empty{}, nil)

		userUCase := NewUseCase(cartClient, productRepo)

		err := userUCase.ChangeProduct(userId, cartArticle)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("ChangeProduct_internal_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)
		cartClient := cart_service.NewMockCartServiceClient(ctrl)

		cartClient.
			EXPECT().
			ChangeProduct(context.Background(), reqCartArticle).
			Return(&empty.Empty{}, errors.ErrInternalError)

		userUCase := NewUseCase(cartClient, productRepo)

		err := userUCase.ChangeProduct(userId, cartArticle)
		assert.Error(t, err, "expected internal error")
	})
}

func TestCartUseCase_DeleteCart(t *testing.T) {
	userId := uint64(2)
	reqUserId := &proto.ReqUserId{
		UserId: userId,
	}

	t.Run("DeleteCart_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)
		cartClient := cart_service.NewMockCartServiceClient(ctrl)

		cartClient.
			EXPECT().
			DeleteCart(context.Background(), reqUserId).
			Return(&empty.Empty{}, nil)

		userUCase := NewUseCase(cartClient, productRepo)

		err := userUCase.DeleteCart(userId)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("DeleteCart_internal_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)
		cartClient := cart_service.NewMockCartServiceClient(ctrl)

		cartClient.
			EXPECT().
			DeleteCart(context.Background(), reqUserId).
			Return(&empty.Empty{}, errors.ErrInternalError)

		userUCase := NewUseCase(cartClient, productRepo)

		err := userUCase.DeleteCart(userId)
		assert.Error(t, err, "expected internal error")
	})
}

func TestCartUseCase_GetPreviewCart(t *testing.T) {
	userId := uint64(2)
	productId := uint64(1)
	reqUserId := &proto.ReqUserId{
		UserId: userId,
	}
	cart := &proto.Cart{
		Products: map[uint64]*proto.ProductPosition{
			productId : &proto.ProductPosition{Count: 1},
		},
	}
	productById := &models.Product{
		Images:       []string{"test"},
	}
	previewCart := &models.PreviewCart{
		Products: []*models.PreviewCartArticle{
			&models.PreviewCartArticle{
				PreviewImage: "test",
				Count:        1,
			},
		},
		Price:    models.TotalPrice{},
	}

	t.Run("GetPreviewCart_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)
		productRepo.
			EXPECT().
			SelectProductById(productId).
			Return(productById, nil)

		cartClient := cart_service.NewMockCartServiceClient(ctrl)
		cartClient.
			EXPECT().
			GetPreviewCart(context.Background(), reqUserId).
			Return(cart, nil)

		userUCase := NewUseCase(cartClient, productRepo)

		res, err := userUCase.GetPreviewCart(userId)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, previewCart, res, "incorrect result")
	})

	t.Run("GetPreviewCart_internal_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)

		cartClient := cart_service.NewMockCartServiceClient(ctrl)
		cartClient.
			EXPECT().
			GetPreviewCart(context.Background(), reqUserId).
			Return(cart, errors.ErrInternalError)

		userUCase := NewUseCase(cartClient, productRepo)

		_, err := userUCase.GetPreviewCart(userId)
		assert.Error(t, err, "expected error")
	})

	t.Run("GetPreviewCart_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)
		productRepo.
			EXPECT().
			SelectProductById(productId).
			Return(productById, errors.ErrInternalError)

		cartClient := cart_service.NewMockCartServiceClient(ctrl)
		cartClient.
			EXPECT().
			GetPreviewCart(context.Background(), reqUserId).
			Return(cart, nil)

		userUCase := NewUseCase(cartClient, productRepo)

		_, err := userUCase.GetPreviewCart(userId)
		assert.Error(t, err, "expected error")
	})
}