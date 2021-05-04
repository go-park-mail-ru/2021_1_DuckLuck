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
		"SELECT id, title, rating, description, base_cost, total_cost, discount, images, id_category "+
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
		&productById.Price.TotalCost,
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
func (r *PostgresqlRepository) GetCountPages(category uint64, count int, filterString string) (int, error) {
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
			"AND c.right_node <= current_node.right_node "+
			filterString+
			" ) ",
		category,
	)

	var countPages int
	if err := row.Scan(&countPages); err != nil {
		return 0, errors.ErrDBInternalError
	}
	countPages = int(math.Ceil(float64(countPages) / float64(count)))

	return countPages, nil
}

// Get count of all pages for this search
func (r *PostgresqlRepository) GetCountSearchPages(category uint64, count int,
	searchString, filterString string) (int, error) {
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
			"AND c.right_node <= current_node.right_node "+
			"AND p.fts @@ plainto_tsquery('ru', $2) "+
			filterString+
			" ) ",
		category,
		searchString,
	)

	var countPages int
	if err := row.Scan(&countPages); err != nil {
		return 0, errors.ErrDBInternalError
	}
	countPages = int(math.Ceil(float64(countPages) / float64(count)))

	return countPages, nil
}

// Create sort string from paginator options
func (r *PostgresqlRepository) CreateSortString(sortKey, sortDirection string) (string, error) {
	// Select order target
	var orderTarget string
	switch sortKey {
	case models.ProductsCostSort:
		orderTarget = "total_cost"
	case models.ProductsRatingSort:
		orderTarget = "rating"
	case models.ProductsDateAddedSort:
		orderTarget = "date_added"
	case models.ProductsDiscountSort:
		orderTarget = "discount"
	default:
		return "", errors.ErrIncorrectPaginator
	}

	// Select order direction
	var orderDirection string
	switch sortDirection {
	case models.PaginatorASC:
		orderDirection = "ASC"
	case models.PaginatorDESC:
		orderDirection = "DESC"
	default:
		return "", errors.ErrIncorrectPaginator
	}

	return fmt.Sprintf("ORDER BY %s %s ", orderTarget, orderDirection), nil
}

// Create filter string from filter options
func (r *PostgresqlRepository) CreateFilterString(filter *models.ProductFilter) string {
	// Check price
	filterString := fmt.Sprintf("AND p.total_cost > %d AND p.total_cost < %d ", filter.MinPrice, filter.MaxPrice)

	// Optional params
	if filter.IsDiscount {
		filterString += "AND p.discount > 0 "
	}
	if filter.IsNew {
		filterString += "AND p.date_added >= date_trunc('month',current_timestamp - interval '1 month') " +
			"AND p.date_added <  date_trunc('month',current_timestamp) "
	}
	if filter.IsRating {
		filterString += "AND p.rating >= 4 "
	}

	return filterString
}

// Select range of products by paginate settings
func (r *PostgresqlRepository) SelectRangeProducts(paginator *models.PaginatorProducts,
	sortString, filterString string) ([]*models.ViewProduct, error) {
	rows, err := r.db.Query(
		"WITH current_node AS ( "+
			"SELECT c.left_node, c.right_node "+
			"FROM categories c "+
			"WHERE c.id = $1 "+
			") "+
			"SELECT p.id, p.title, p.base_cost, p.total_cost, p.discount, p.rating, p.images[1] "+
			"FROM current_node, products p "+
			"JOIN categories c ON c.id = p.id_category "+
			"WHERE (c.left_node >= current_node.left_node "+
			"AND c.right_node <= current_node.right_node "+
			filterString+
			" ) "+
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
			&product.Price.TotalCost,
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

// Find list of products by query string
func (r *PostgresqlRepository) SearchRangeProducts(searchQuery *models.SearchQuery,
	sortString, filterString string) ([]*models.ViewProduct, error) {
	rows, err := r.db.Query(
		"WITH current_node AS ( "+
			"SELECT c.left_node, c.right_node "+
			"FROM categories c "+
			"WHERE c.id = $1 "+
			") "+
			"SELECT p.id, p.title, p.base_cost, p.total_cost, p.discount, p.rating, p.images[1] "+
			"FROM current_node, products p "+
			"JOIN categories c ON c.id = p.id_category "+
			"WHERE (c.left_node >= current_node.left_node "+
			"AND c.right_node <= current_node.right_node "+
			"AND p.fts @@ plainto_tsquery('ru', $2) "+
			filterString+
			" ) "+
			sortString+
			"LIMIT $3 OFFSET $4",
		searchQuery.Category,
		searchQuery.QueryString,
		searchQuery.Count,
		searchQuery.Count*(searchQuery.PageNum-1),
	)
	if err != nil {
		return nil, errors.ErrIncorrectSearchQuery
	}
	defer rows.Close()

	products := make([]*models.ViewProduct, 0)
	for rows.Next() {
		product := &models.ViewProduct{}
		err = rows.Scan(
			&product.Id,
			&product.Title,
			&product.Price.BaseCost,
			&product.Price.TotalCost,
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
