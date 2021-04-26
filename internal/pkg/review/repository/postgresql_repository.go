package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/review"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) review.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

// Select all reviews by product id
func (r *PostgresqlRepository) SelectReviewsByProductId(productId uint64) ([]*models.ViewReview, error) {

}

// Add new review for product
func (r *PostgresqlRepository) AddReviewProduct(userId uint64, review *models.Review) error {

}
