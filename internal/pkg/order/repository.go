package order

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order Repository

type Repository interface {
	AddOrder(order *models.Order, userId uint64, products []*models.PreviewCartArticle,
		price *models.TotalPrice) (*models.OrderNumber, error)
	SelectRangeOrders(orderId uint64, sortString string, paginator *models.PaginatorOrders) ([]*models.PlacedOrder, error)
	CreateSortString(sortKey, sortDirection string) (string, error)
	GetCountPages(userId uint64, countOrdersOnPage int) (int, error)
	GetProductsInOrder(orderId uint64) ([]*models.PreviewOrderedProducts, error)
	ChangeStatusOrder(orderId uint64, status string) (*models.OrderNumber, uint64, error)
}
