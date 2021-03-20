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

func (pr *PostgresqlRepository) GetNextLevelCategories(categoryId uint64) ([]*models.CategoriesCatalog, error) {
	rows, err := pr.db.Query(
		"SELECT c.id, c.name "+
			"FROM category c "+
			"JOIN subsetCategory s1 ON c.id = s1.idCategory "+
			"JOIN subsetCategory s2 on s1.idCategory = s2.idCategory "+
			"WHERE (s1.idSubSet = $1) and (s2.idCategory = s2.idSubset) "+
			"and (s2.level = s1.level + 1)"+
			"GROUP BY c.id, c.name "+
			"ORDER BY c.name",
		categoryId,
	)
	defer rows.Close()

	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	categories := make([]*models.CategoriesCatalog, 0)
	for rows.Next() {
		nextLevelCategory := &models.CategoriesCatalog{}
		err = rows.Scan(
			&nextLevelCategory.Id,
			&nextLevelCategory.LevelName,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, nextLevelCategory)
	}

	return categories, nil
}

func (pr *PostgresqlRepository) GetCategoriesByLevel(level uint64) ([]*models.CategoriesCatalog, error) {
	rows, err := pr.db.Query(
		"SELECT c.id, c.name "+
			"FROM category c "+
			"JOIN subsetCategory s1 ON c.id = s1.idCategory "+
			"GROUP BY c.id, c.name "+
			"HAVING count(*) = $1",
		level,
	)
	defer rows.Close()

	if err != nil {
		return nil, errors.ErrDBInternalError
	}

	categories := make([]*models.CategoriesCatalog, 0)
	for rows.Next() {
		nextLevelCategory := &models.CategoriesCatalog{}
		err = rows.Scan(
			&nextLevelCategory.Id,
			&nextLevelCategory.LevelName,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, nextLevelCategory)
	}

	return categories, nil
}
