package repository

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/review"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) review.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

// Select range of reviews
func (r *PostgresqlRepository) SelectRangeReviews(productId uint64, sortString string,
	paginator *models.PaginatorReviews) ([]*models.ViewReview, error) {
	rows, err := r.db.Query(
		"SELECT rating, advantages, disadvantages, comment, is_public, "+
			"date_added, user_id "+
			"FROM reviews "+
			"WHERE product_id = $1 "+
			sortString+
			"LIMIT $2 OFFSET $3",
		productId,
		paginator.Count,
		paginator.Count*(paginator.PageNum-1),
	)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}
	defer rows.Close()

	reviews := make([]*models.ViewReview, 0)
	for rows.Next() {
		userReview := &models.ViewReview{}
		err = rows.Scan(
			&userReview.Rating,
			&userReview.Advantages,
			&userReview.Disadvantages,
			&userReview.Comment,
			&userReview.IsPublic,
			&userReview.DateAdded,
			&userReview.UserId,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, userReview)
	}

	return reviews, nil
}

// Get count of all review pages for this product
func (r *PostgresqlRepository) GetCountPages(productId uint64, countOrdersOnPage int) (int, error) {
	row := r.db.QueryRow(
		"SELECT count(id) "+
			"FROM reviews "+
			"WHERE product_id = $1",
		productId,
	)

	var countPages int
	if err := row.Scan(&countPages); err != nil {
		return 0, errors.ErrDBInternalError
	}
	countPages = int(math.Ceil(float64(countPages) / float64(countOrdersOnPage)))

	return countPages, nil
}

// Create sort string for query
func (r *PostgresqlRepository) CreateSortString(sortKey, sortDirection string) (string, error) {
	// Select order target
	var orderTarget string
	switch sortKey {
	case models.ReviewDateAddedSort:
		orderTarget = "date_added"
	default:
		return "", errors.ErrIncorrectPaginator
	}

	// Select order direction
	var orderDirection string
	switch sortDirection {
	case models.ReviewPaginatorASC:
		orderDirection = "ASC"
	case models.ReviewPaginatorDESC:
		orderDirection = "DESC"
	default:
		return "", errors.ErrIncorrectPaginator
	}

	return fmt.Sprintf("ORDER BY %s %s ", orderTarget, orderDirection), nil
}

// Select all statistics about reviews by product id
func (r *PostgresqlRepository) SelectStatisticsByProductId(productId uint64) (*models.ReviewStatistics, error) {
	rows, err := r.db.Query(
		"SELECT count(id), rating "+
			"FROM reviews "+
			"WHERE product_id = $1 "+
			"GROUP BY rating "+
			"ORDER BY rating",
		productId,
	)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}
	defer rows.Close()

	statistics := &models.ReviewStatistics{}
	statistics.Stars = make([]int, 5, 5)
	var countStars int
	var rating int
	for rows.Next() {
		err = rows.Scan(
			&countStars,
			&rating,
		)
		if err != nil {
			return nil, err
		}
		statistics.Stars[rating] = countStars
	}

	return statistics, nil
}

// Check rights for review (the user has completed orders)
func (r *PostgresqlRepository) CheckReview(userId uint64, productId uint64) bool {
	row := r.db.QueryRow(
		"SELECT count(us.id) "+
			"FROM user_orders us "+
			"JOIN ordered_products op ON us.id = op.order_id " +
			"LEFT JOIN reviews rv ON op.product_id = rv.product_id " +
			"WHERE (us.user_id = $1 AND op.product_id = $2 " +
			"AND rv.product_id IS NULL)",
		userId,
		productId,
	)

	var isExist int
	if err := row.Scan(&isExist); err != nil || isExist == 0 {
		return false
	}

	return true
}

// Add new review for product
func (r *PostgresqlRepository) AddReview(userId uint64, review *models.Review) (uint64, error) {
	row := r.db.QueryRow(
		"INSERT INTO reviews(product_id, user_id, rating, advantages, "+
			"disadvantages, comment, is_public) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		review.ProductId,
		userId,
		review.Rating,
		review.Advantages,
		review.Disadvantages,
		review.Comment,
		review.IsPublic,
	)

	var reviewId uint64
	if err := row.Scan(&reviewId); err != nil {
		return 0, errors.ErrDBInternalError
	}

	return reviewId, nil
}
