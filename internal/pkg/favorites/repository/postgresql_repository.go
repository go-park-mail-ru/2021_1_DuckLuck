package repository

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/favorites"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) favorites.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

func (r *PostgresqlRepository) AddProductToFavorites(productId, userId uint64) error {
	_, err := r.db.Exec(
		"INSERT INTO favorites(product_id, user_id) "+
			"VALUES ($1, $2)",
		productId,
		userId,
	)
	if err != nil {
		return errors.ErrDBInternalError
	}

	return nil
}

func (r *PostgresqlRepository) DeleteProductFromFavorites(productId, userId uint64) error {
	_, err := r.db.Exec(
		"DELETE FROM favorites "+
			"WHERE  (product_id = $1 AND user_id = $2)",
		productId,
		userId,
	)
	if err != nil {
		return errors.ErrDBInternalError
	}

	return nil
}

func (r *PostgresqlRepository) GetCountPages(userId uint64, count int) (int, error) {
	row := r.db.QueryRow(
		"SELECT count(*) "+
			"FROM favorites "+
			"WHERE user_id = $1",
		userId,
	)

	var countPages int
	if err := row.Scan(&countPages); err != nil {
		return 0, errors.ErrDBInternalError
	}
	countPages = int(math.Ceil(float64(countPages) / float64(count)))

	return countPages, nil
}

func (r *PostgresqlRepository) CreateSortString(sortKey, sortDirection string) (string, error) {
	// Select order target
	var orderTarget string
	switch sortKey {
	case models.FavoritesCostSort:
		orderTarget = "total_cost"
	case models.FavoritesRatingSort:
		orderTarget = "(CASE WHEN avg_rating IS NULL THEN 0 ELSE avg_rating END)"
	case models.FavoritesDateAddedSort:
		orderTarget = "date_added"
	case models.FavoritesDiscountSort:
		orderTarget = "discount"
	default:
		return "", errors.ErrIncorrectPaginator
	}

	// Select order direction
	var orderDirection string
	switch sortDirection {
	case models.FavoritesPaginatorASC:
		orderDirection = "ASC"
	case models.FavoritesPaginatorDESC:
		orderDirection = "DESC"
	default:
		return "", errors.ErrIncorrectPaginator
	}

	return fmt.Sprintf("ORDER BY %s %s ", orderTarget, orderDirection), nil
}

func (r *PostgresqlRepository) SelectRangeFavorites(paginator *models.PaginatorFavorites,
	sortString string, userId uint64) ([]*models.ViewFavorite, error) {
	rows, err := r.db.Query(
		"SELECT p.id, p.title, p.base_cost, p.total_cost, "+
			"p.discount, p.images[1], "+
			"avg_rating, count_reviews "+
			"FROM products p "+
			"JOIN categories c ON c.id = p.id_category "+
			"JOIN favorites f ON (f.product_id = p.id AND f.user_id = $1) "+
			"LEFT JOIN ( "+
			"	SELECT product_id, "+
			"	AVG(rating) as avg_rating, "+
			"	COUNT(*) as count_reviews "+
			"FROM reviews "+
			"GROUP BY product_id "+
			") AS R ON P.id = R.product_id "+
			sortString+
			"LIMIT $2 OFFSET $3",
		userId,
		paginator.Count,
		paginator.Count*(paginator.PageNum-1),
	)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}
	defer rows.Close()

	products := make([]*models.ViewFavorite, 0)
	rating := sql.NullFloat64{}
	countReviews := sql.NullInt64{}
	for rows.Next() {
		product := &models.ViewFavorite{}
		err = rows.Scan(
			&product.Id,
			&product.Title,
			&product.Price.BaseCost,
			&product.Price.TotalCost,
			&product.Price.Discount,
			&product.PreviewImage,
			&rating,
			&countReviews,
		)
		product.Rating = float32(rating.Float64)
		product.CountReviews = uint64(countReviews.Int64)

		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *PostgresqlRepository) GetUserFavorites(userId uint64) (*models.UserFavorites, error) {
	rows, err := r.db.Query(
		"SELECT product_id "+
			"FROM favorites "+
			"WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}
	defer rows.Close()

	favoritesProducts := make([]uint64, 0)
	var productId uint64
	for rows.Next() {
		err = rows.Scan(
			&productId,
		)

		if err != nil {
			return nil, err
		}
		favoritesProducts = append(favoritesProducts, productId)
	}

	return &models.UserFavorites{
		Products: favoritesProducts,
	}, nil
}
