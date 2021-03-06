package category

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category Repository

type Repository interface {
	GetNextLevelCategories(categoryId uint64) ([]*models.CategoriesCatalog, error)
	GetCategoriesByLevel(level uint64) ([]*models.CategoriesCatalog, error)
	GetBordersOfBranch(categoryId uint64) (uint64, uint64, error)
	GetPathToCategory(categoryId uint64) ([]*models.CategoriesCatalog, error)
}
