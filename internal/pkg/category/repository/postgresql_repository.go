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
		"SELECT c.id, c.name "+
			"FROM category c "+
			"JOIN subset_category s1 ON c.id = s1.id_category "+
			"JOIN subset_category s2 on s1.id_category = s2.id_category "+
			"WHERE (s1.id_subset = $1) and (s2.id_category = s2.id_subset) "+
			"and (s2.level = s1.level + 1)"+
			"GROUP BY c.id, c.name "+
			"ORDER BY c.name",
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
			"FROM category c "+
			"JOIN subset_category s1 ON c.id = s1.id_category "+
			"GROUP BY c.id, c.name "+
			"HAVING count(*) = $1",
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
		"SELECT id_category "+
			"FROM subsets_category "+
			"WHERE id_subset = $1",
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
		"SELECT c.id, c.name FROM  subsets_category s "+
			"LEFT JOIN category c ON c.id = s.id_subset "+
			"WHERE s.id_category = $1 "+
			"ORDER BY s.level",
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
