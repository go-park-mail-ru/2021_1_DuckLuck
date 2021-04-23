package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) category.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

// Get lower level in categories tree
func (r *PostgresqlRepository) GetNextLevelCategories(categoryId uint64) ([]*models.CategoriesCatalog, error) {
	rows, err := r.db.Query(
		"WITH current_node AS ( "+
			"SELECT c.left_node, c.right_node, c.level + 1  as level "+
			"FROM categories c "+
			"WHERE c.id = $1 "+
			") "+
			"SELECT c.id, c.name "+
			"FROM categories c, current_node "+
			"WHERE (c.left_node > current_node.left_node "+
			"AND c.right_node < current_node.right_node "+
			"AND c.level = current_node.level)",
		categoryId,
	)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}
	defer rows.Close()

	categories := make([]*models.CategoriesCatalog, 0)
	for rows.Next() {
		nextLevelCategory := &models.CategoriesCatalog{}
		err = rows.Scan(
			&nextLevelCategory.Id,
			&nextLevelCategory.Name,
		)
		if err != nil {
			return nil, errors.ErrDBInternalError
		}
		categories = append(categories, nextLevelCategory)
	}

	return categories, nil
}

// Get categories in select level
func (r *PostgresqlRepository) GetCategoriesByLevel(level uint64) ([]*models.CategoriesCatalog, error) {
	rows, err := r.db.Query(
		"SELECT c.id, c.name "+
			"FROM categories c "+
			"WHERE c.level = $1",
		level,
	)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}
	defer rows.Close()

	categories := make([]*models.CategoriesCatalog, 0)
	for rows.Next() {
		nextLevelCategory := &models.CategoriesCatalog{}
		err = rows.Scan(
			&nextLevelCategory.Id,
			&nextLevelCategory.Name,
		)
		if err != nil {
			return nil, errors.ErrDBInternalError
		}
		categories = append(categories, nextLevelCategory)
	}

	return categories, nil
}

// Get left and right border of branch
func (r *PostgresqlRepository) GetBordersOfBranch(categoryId uint64) (uint64, uint64, error) {
	row := r.db.QueryRow(
		"SELECT c.left_node, c.right_node "+
			"FROM categories c "+
			"WHERE c.id = $1",
		categoryId,
	)

	var left, right uint64
	err := row.Scan(
		&left,
		&right,
	)

	if err != nil {
		return 0, 0, errors.ErrDBInternalError
	}

	return left, right, nil
}

// Get path from root to category
func (r *PostgresqlRepository) GetPathToCategory(categoryId uint64) ([]*models.CategoriesCatalog, error) {
	rows, err := r.db.Query(
		"WITH current_node AS ( "+
			"SELECT c.left_node, c.right_node, c.level + 1  as level "+
			"FROM categories c "+
			"WHERE c.id = $1 "+
			") "+
			"SELECT c.id, c.name "+
			"FROM categories c, current_node "+
			"WHERE (c.left_node <= current_node.left_node "+
			"AND c.right_node >= current_node.right_node "+
			"AND c.level BETWEEN 1 AND current_node.level)",
		categoryId,
	)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}
	defer rows.Close()

	categories := make([]*models.CategoriesCatalog, 0)
	for rows.Next() {
		nextLevelCategory := &models.CategoriesCatalog{}
		err = rows.Scan(
			&nextLevelCategory.Id,
			&nextLevelCategory.Name,
		)
		if err != nil {
			return nil, errors.ErrDBInternalError
		}
		categories = append(categories, nextLevelCategory)
	}

	return categories, nil
}
