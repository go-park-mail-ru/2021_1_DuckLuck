package order

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order UseCase

type UseCase interface {
	GetPreviewOrder(userId uint64, previewCart *models.PreviewCart) (*models.PreviewOrder, error)
	AddCompletedOrder(order *models.Order, userId uint64, previewCart *models.PreviewCart) (*models.OrderNumber, error)
	GetRangeOrders(userId uint64, paginator *models.PaginatorOrders) (*models.RangeOrders, error)
}
