package repository

import (
	"database/sql"
	"fmt"
	"math"

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

// Get count of all pages for this category
func (r *PostgresqlRepository) GetCountPages(paginator *models.PaginatorProducts) (int, error) {
	row := r.db.QueryRow(
		"WITH current_node AS ( "+
			"SELECT c.left_node, c.right_node "+
			"FROM categories c "+
			"WHERE c.id = $1 "+
			") "+
			"SELECT count(p.id) "+
			"FROM current_node, products p "+
			"JOIN categories c ON c.id = p.id_category "+
			"WHERE (c.left_node >= current_node.left_node "+
			"AND c.right_node <= current_node.right_node)",
		paginator.Category,
	)

	var countPages int
	if err := row.Scan(&countPages); err != nil {
		return 0, errors.ErrDBInternalError
	}
	countPages = int(math.Ceil(float64(countPages) / float64(paginator.Count)))

	return countPages, nil
}

// Create sort string from paginator options
func (r *PostgresqlRepository) CreateSortString(paginator *models.PaginatorProducts) (string, error) {
	// Select order target
	var orderTarget string
	switch paginator.SortKey {
	case models.ProductsCostSort:
		orderTarget = "base_cost"
	case models.ProductsRatingSort:
		orderTarget = "rating"
	default:
		return "", errors.ErrIncorrectPaginator
	}

	// Select order direction
	var orderDirection string
	switch paginator.SortDirection {
	case models.PaginatorASC:
		orderDirection = "ASC"
	case models.PaginatorDESC:
		orderDirection = "DESC"
	default:
		return "", errors.ErrIncorrectPaginator
	}

	return fmt.Sprintf("ORDER BY %s %s ", orderTarget, orderDirection), nil
}

// Select range of products by paginate settings
func (r *PostgresqlRepository) SelectRangeProducts(paginator *models.PaginatorProducts,
	sortString string) ([]*models.ViewProduct, error) {
	rows, err := r.db.Query(
		"WITH current_node AS ( "+
			"SELECT c.left_node, c.right_node "+
			"FROM categories c "+
			"WHERE c.id = $1 "+
			") "+
			"SELECT p.id, p.title, p.base_cost, p.discount, p.rating, p.images[1] "+
			"FROM current_node, products p "+
			"JOIN categories c ON c.id = p.id_category "+
			"WHERE (c.left_node >= current_node.left_node "+
			"AND c.right_node <= current_node.right_node) "+
			sortString+
			"LIMIT $2 OFFSET $3",
		paginator.Category,
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

	return products, nil
}
