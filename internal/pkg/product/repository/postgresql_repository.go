package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/lib/pq"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) product.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

// Select one product by id
func (r *PostgresqlRepository) SelectProductById(productId uint64) (*models.Product, error) {
	row := r.db.QueryRow(
		"SELECT id, title, rating, description, base_cost, discount, images, id_category "+
			"FROM products WHERE id = $1",
		productId,
	)

	description := sql.NullString{}
	productById := models.Product{}
	err := row.Scan(
		&productById.Id,
		&productById.Title,
		&productById.Rating,
		&description,
		&productById.Price.BaseCost,
		&productById.Price.Discount,
		pq.Array(&productById.Images),
		&productById.Category,
	)
	productById.Description = description.String

	if err != nil {
		return nil, errors.ErrDBInternalError
	}

	return &productById, nil
}

// Select range of products by paginate settings
func (r *PostgresqlRepository) SelectRangeProducts(paginator *models.PaginatorProducts,
	categories *[]uint64) (*models.RangeProducts, error) {
	row := r.db.QueryRow(
		"SELECT ceil(count(*) / $1) FROM products "+
			"WHERE id_category = ANY($2)",
		paginator.Count,
		pq.Array(*categories),
	)

	var countPages int
	if err := row.Scan(&countPages); err != nil {
		return nil, errors.ErrDBInternalError
	}

	var sortString string
	switch paginator.SortKey {
	case models.ProductsCostSort:
		if paginator.SortDirection == models.PaginatorASC {
			sortString = fmt.Sprintf("ORDER BY base_cost ASC ")
		} else if paginator.SortDirection == models.PaginatorDESC {
			sortString = fmt.Sprintf("ORDER BY base_cost DESC ")
		}
	case models.ProductsRatingSort:
		if paginator.SortDirection == models.PaginatorASC {
			sortString = fmt.Sprintf("ORDER BY rating ASC ")
		} else if paginator.SortDirection == models.PaginatorDESC {
			sortString = fmt.Sprintf("ORDER BY rating DESC ")
		}
	}

	rows, err := r.db.Query(
		"SELECT id, title, base_cost, discount, rating, images[1] "+
			"FROM products "+
			"WHERE id_category = ANY($1) "+
			sortString+
			"LIMIT $2 OFFSET $3",
		pq.Array(*categories),
		paginator.Count,
		paginator.Count*(paginator.PageNum-1),
	)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}
	defer rows.Close()

	products := make([]*models.ViewProduct, 0)
	for rows.Next() {
		product := &models.ViewProduct{}
		err = rows.Scan(
			&product.Id,
			&product.Title,
			&product.Price.BaseCost,
			&product.Price.Discount,
			&product.Rating,
			&product.PreviewImage,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return &models.RangeProducts{
		ListPreviewProducts: products,
		MaxCountPages:       countPages,
	}, nil
}
