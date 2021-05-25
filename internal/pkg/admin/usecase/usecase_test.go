package usecase

import (
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	notification_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/notification/mock"
	order_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCategoryUseCase_ChangeOrderStatus(t *testing.T) {
	updateOrder := models.UpdateOrder{
		OrderId: 3,
		Status:  "получено",
	}
	orderNumber := models.OrderNumber{Number: "0001-00000001"}
	userId := uint64(2)
	subscribes := models.Subscribes{
		Credentials: nil,
	}

	t.Run("ChangeOrderStatus_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orderRepo := order_mock.NewMockRepository(ctrl)
		orderRepo.
			EXPECT().
			ChangeStatusOrder(updateOrder.OrderId, updateOrder.Status).
			Return(&orderNumber, userId, nil)

		notificationRepo := notification_mock.NewMockRepository(ctrl)
		notificationRepo.
			EXPECT().
			SelectCredentialsByUserId(userId).
			Return(&subscribes, nil)

		adminUCase := NewUseCase(notificationRepo, orderRepo)

		err := adminUCase.ChangeOrderStatus(&updateOrder)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("ChangeOrderStatus_can't_change_status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orderRepo := order_mock.NewMockRepository(ctrl)
		orderRepo.
			EXPECT().
			ChangeStatusOrder(updateOrder.OrderId, updateOrder.Status).
			Return(&orderNumber, userId, errors.ErrInternalError)

		notificationRepo := notification_mock.NewMockRepository(ctrl)

		adminUCase := NewUseCase(notificationRepo, orderRepo)

		err := adminUCase.ChangeOrderStatus(&updateOrder)
		assert.Error(t, err, "unexpected error")
	})
}
