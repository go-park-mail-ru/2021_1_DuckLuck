package repository

import (
	"database/sql"
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

func (r *PostgresqlRepository) GetOrders(userId uint64) ([]*models.Order, error) {
	rows, err := r.db.Query(
		"SELECT first_name, last_name, email, address "+
			"FROM user_orders "+
			"WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}
	defer rows.Close()

	orders := make([]*models.Order, 0)
	for rows.Next() {
		order := &models.Order{}
		err = rows.Scan(
			&order.Recipient.FirstName,
			&order.Recipient.LastName,
			&order.Recipient.Email,
			&order.Address.Address,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Add order in db
func (r *PostgresqlRepository) AddOrder(order *models.Order, userId uint64,
	products []*models.PreviewCartArticle, price *models.TotalPrice) (uint64, error) {
	row := r.db.QueryRow(
		"INSERT INTO user_orders(user_id, first_name, last_name, email, "+
			"address, base_cost, total_cost, discount) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		userId,
		order.Recipient.FirstName,
		order.Recipient.LastName,
		order.Recipient.Email,
		order.Address.Address,
		price.TotalBaseCost,
		price.TotalCost,
		price.TotalDiscount,
	)
	var orderId uint64
	if err := row.Scan(&orderId); err != nil {
		return 0, errors.ErrDBInternalError
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
			return 0, errors.ErrDBInternalError
		}
	}

	return orderId, nil
}
