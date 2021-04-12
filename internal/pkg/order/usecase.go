package order

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type UseCase interface {
	GetPreviewOrder(userId uint64, previewCart *models.PreviewCart) (*models.PreviewOrder, error)
	AddCompletedOrder(order *models.Order, userId uint64, previewCart *models.PreviewCart) (uint64, error)
}
