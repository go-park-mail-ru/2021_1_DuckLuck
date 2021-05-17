package repository

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) order.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

func (r *PostgresqlRepository) SelectRangeOrders(userId uint64, sortString string,
	paginator *models.PaginatorOrders) ([]*models.PlacedOrder, error) {
	rows, err := r.db.Query(
		"SELECT id, address, total_cost, date_added, date_delivery, "+
			"order_num, status_pay, status_delivery "+
			"FROM user_orders "+
			"WHERE user_id = $1 "+
			sortString+
			"LIMIT $2 OFFSET $3",
		userId,
		paginator.Count,
		paginator.Count*(paginator.PageNum-1),
	)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}
	defer rows.Close()

	orders := make([]*models.PlacedOrder, 0)
	for rows.Next() {
		placedOrder := &models.PlacedOrder{}
		err = rows.Scan(
			&placedOrder.Id,
			&placedOrder.Address.Address,
			&placedOrder.TotalCost,
			&placedOrder.DateAdded,
			&placedOrder.DateDelivery,
			&placedOrder.OrderNumber.Number,
			&placedOrder.StatusPay,
			&placedOrder.StatusDelivery,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, placedOrder)
	}

	return orders, nil
}

// Get count of all pages for this category
func (r *PostgresqlRepository) GetCountPages(userId uint64, countOrdersOnPage int) (int, error) {
	row := r.db.QueryRow(
		"SELECT count(id) "+
			"FROM user_orders "+
			"WHERE user_id = $1",
		userId,
	)

	var countPages int
	if err := row.Scan(&countPages); err != nil {
		return 0, errors.ErrDBInternalError
	}
	countPages = int(math.Ceil(float64(countPages) / float64(countOrdersOnPage)))

	return countPages, nil
}

func (r *PostgresqlRepository) CreateSortString(sortKey, sortDirection string) (string, error) {
	// Select order target
	var orderTarget string
	switch sortKey {
	case models.OrdersDateAddedSort:
		orderTarget = "date_added"
	default:
		return "", errors.ErrIncorrectPaginator
	}

	// Select order direction
	var orderDirection string
	switch sortDirection {
	case models.OrdersPaginatorASC:
		orderDirection = "ASC"
	case models.OrdersPaginatorDESC:
		orderDirection = "DESC"
	default:
		return "", errors.ErrIncorrectPaginator
	}

	return fmt.Sprintf("ORDER BY %s %s ", orderTarget, orderDirection), nil
}

func (r *PostgresqlRepository) GetProductsInOrder(orderId uint64) ([]*models.PreviewOrderedProducts, error) {
	rows, err := r.db.Query(
		"SELECT p.id, p.images[1] "+
			"FROM ordered_products rp "+
			"JOIN products p ON rp.product_id = p.id "+
			"WHERE rp.order_id = $1",
		orderId,
	)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}
	defer rows.Close()

	products := make([]*models.PreviewOrderedProducts, 0)
	for rows.Next() {
		orderedProduct := &models.PreviewOrderedProducts{}
		err = rows.Scan(
			&orderedProduct.Id,
			&orderedProduct.PreviewImage,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, orderedProduct)
	}

	return products, nil
}

// Add order in db
func (r *PostgresqlRepository) AddOrder(order *models.Order, userId uint64,
	products []*models.PreviewCartArticle, price *models.TotalPrice) (*models.OrderNumber, error) {
	row := r.db.QueryRow(
		"INSERT INTO user_orders(user_id, first_name, last_name, email, "+
			"address, base_cost, total_cost, discount) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, order_num",
		userId,
		order.Recipient.FirstName,
		order.Recipient.LastName,
		order.Recipient.Email,
		order.Address.Address,
		price.TotalBaseCost,
		price.TotalCost,
		price.TotalDiscount,
	)
	var orderNumber models.OrderNumber
	var orderId int
	if err := row.Scan(&orderId, &orderNumber.Number); err != nil {
		return nil, errors.ErrDBInternalError
	}

	for _, item := range products {
		res := r.db.QueryRow(
			"INSERT INTO ordered_products(product_id, order_id, num, base_cost, discount) "+
				"VALUES ($1, $2, $3, $4, $5) RETURNING id",
			item.Id,
			orderId,
			item.Count,
			item.Price.BaseCost,
			item.Price.Discount,
		)
		if res.Err() != nil {
			return nil, errors.ErrDBInternalError
		}
	}

	return &orderNumber, nil
}

func (r *PostgresqlRepository) ChangeStatusOrder(orderId uint64,
	status string) (*models.OrderNumber, uint64, error) {
	row := r.db.QueryRow(
		"UPDATE user_orders " +
			"SET status = $1 "+
			"WHERE id = $2 " +
			"RETURNING user_id, order_num",
			status,
			orderId,
	)

	var orderNumber models.OrderNumber
	var userId uint64
	if err := row.Scan(&userId, &orderNumber.Number); err != nil {
		return nil, 0, errors.ErrDBInternalError
	}

	return &orderNumber, userId, nil
}
