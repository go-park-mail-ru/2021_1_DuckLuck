package order

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type Repository interface {
	AddOrder(order *models.Order, userId uint64, products []*models.PreviewCartArticle,
		price *models.TotalPrice) (uint64, error)
}
