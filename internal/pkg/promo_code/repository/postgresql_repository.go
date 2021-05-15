package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/promo_code"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) promo_code.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

func (r *PostgresqlRepository) GetDiscountPriceByPromo(productId uint64, promoCode string) (*models.ProductPrice, error) {
	row := r.db.QueryRow(
		"WITH pr AS ( "+
			"    SELECT id, sale "+
			"    FROM promo_codes "+
			"    WHERE code = $2 "+
			") "+
			"SELECT p.base_cost, p.discount, pr.sale "+
			"FROM products p "+
			"JOIN pr ON (pr.id = ANY(p.sale_group) AND p.id = $1)",
		productId,
		promoCode,
	)

	promoSale := sql.NullInt64{}
	var baseCost, discount int
	err := row.Scan(
		&baseCost,
		&discount,
		&promoSale,
	)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	discount = int(float32(discount) * (1.0 + float32(promoSale.Int64)) / 100.0)
	totalCost := baseCost - discount
	if discount > baseCost {
		discount = baseCost
		totalCost = 0
	}

	return &models.ProductPrice{
		Discount:  discount,
		BaseCost:  baseCost,
		TotalCost: totalCost,
	}, nil
}
