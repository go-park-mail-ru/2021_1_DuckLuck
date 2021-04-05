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
		"SELECT id, title, rating, description, baseCost, discount, images, idCategory "+
			"FROM products WHERE id = $1",
		productId,
	)

	var idCategory int
	productById := models.Product{}
	err := row.Scan(
		&productById.Id,
		&productById.Title,
		&productById.Rating,
		&productById.Description,
		&productById.Price.BaseCost,
		&productById.Price.Discount,
		pq.Array(&productById.Images),
		&idCategory,
	)

	rows, err := r.db.Query(
		"SELECT c.id, c.name FROM  subsetCategory s "+
			"LEFT JOIN category c ON c.id = s.idSubset "+
			"WHERE s.idCategory = $1 "+
			"ORDER BY s.level",
		idCategory,
	)
	defer rows.Close()

	if err != nil {
		return nil, errors.ErrDBInternalError
	}

	for rows.Next() {
		nextLevelCategory := &models.CategoriesCatalog{}
		err = rows.Scan(
			&nextLevelCategory.Id,
			&nextLevelCategory.Name,
		)
		if err != nil {
			return nil, err
		}
		productById.Category = append(productById.Category, nextLevelCategory)
	}

	return &productById, nil
}

// Select range of products by paginate settings
func (r *PostgresqlRepository) SelectRangeProducts(paginator *models.PaginatorProducts) (*models.RangeProducts, error) {
	row := r.db.QueryRow(
		"SELECT ceil(count(*) / $1) FROM products",
		paginator.Count,
	)

	var countPages int
	if err := row.Scan(&countPages); err != nil {
		return nil, errors.ErrDBInternalError
	}

	var sortString string
	switch paginator.SortKey {
	case models.ProductsCostSort:
		if paginator.SortDirection == models.PaginatorASC {
			sortString = fmt.Sprintf("ORDER BY baseCost ASC ")
		} else if paginator.SortDirection == models.PaginatorDESC {
			sortString = fmt.Sprintf("ORDER BY baseCost DESC ")
		}
	case models.ProductsRatingSort:
		if paginator.SortDirection == models.PaginatorASC {
			sortString = fmt.Sprintf("ORDER BY rating ASC ")
		} else if paginator.SortDirection == models.PaginatorDESC {
			sortString = fmt.Sprintf("ORDER BY rating DESC ")
		}
	}

	rows, err := r.db.Query(
		"SELECT id, title, baseCost, discount, rating, images[1] "+
			"FROM products "+
			sortString+
			"LIMIT $1 OFFSET $2",
		paginator.Count,
		paginator.Count*(paginator.PageNum-1),
	)
	defer rows.Close()

	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

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
