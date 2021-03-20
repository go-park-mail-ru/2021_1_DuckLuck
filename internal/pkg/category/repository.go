package category

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type Repository interface {
	GetNextLevelCategories(categoryId uint64) ([]*models.CategoriesCatalog, error)
	GetCategoriesByLevel(level uint64) ([]*models.CategoriesCatalog, error)
}
