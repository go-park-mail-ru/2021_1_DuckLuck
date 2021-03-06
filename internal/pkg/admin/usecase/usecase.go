package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/admin"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/notification"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/server_push"
)

type AdminUseCase struct {
	NotificationRepo notification.Repository
	OrderRepo        order.Repository
}

func NewUseCase(notificationRepo notification.Repository, orderRepo order.Repository) admin.UseCase {
	return &AdminUseCase{
		NotificationRepo: notificationRepo,
		OrderRepo:        orderRepo,
	}
}

func (u *AdminUseCase) ChangeOrderStatus(updateOrder *models.UpdateOrder) error {
	orderNumber, userId, err := u.OrderRepo.ChangeStatusOrder(updateOrder.OrderId, updateOrder.Status)
	if err != nil {
		return errors.ErrInternalError
	}

	subscribes, err := u.NotificationRepo.SelectCredentialsByUserId(userId)
	if err == nil {
		for endpoint, keys := range subscribes.Credentials {
			err = server_push.Push(&server_push.Subscription{
				Endpoint: endpoint,
				Auth:     keys.Auth,
				P256dh:   keys.P256dh,
			}, models.OrderNotification{
				Number: *orderNumber,
				Status: updateOrder.Status,
			})

			if err != nil {
				return errors.ErrInternalError
			}
		}
	}

	return nil
}
