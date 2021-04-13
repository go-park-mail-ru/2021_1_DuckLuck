package order

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order Repository

type Repository interface {
	AddOrder(order *models.Order, userId uint64, products []*models.PreviewCartArticle,
		price *models.TotalPrice) (uint64, error)
}
