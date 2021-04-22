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

// Get id of all subcategories
func (r *PostgresqlRepository) GetAllSubCategoriesId(categoryId uint64) ([]uint64, error) {
	rows, err := r.db.Query(
		"WITH current_node AS ( "+
			"SELECT c.left_node, c.right_node "+
			"FROM categories c "+
			"WHERE c.id = $1 "+
			") "+
			"SELECT c.id "+
			"FROM categories c, current_node "+
			"WHERE (c.left_node > current_node.left_node "+
			"AND c.right_node < current_node.right_node)",
		categoryId,
	)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}
	defer rows.Close()

	categoriesId := make([]uint64, 0)
	for rows.Next() {
		var id uint64
		err = rows.Scan(
			&id,
		)
		if err != nil {
			return nil, errors.ErrDBInternalError
		}
		categoriesId = append(categoriesId, id)
	}

	return categoriesId, nil
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
			"AND c.level <= current_node.level)",
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
