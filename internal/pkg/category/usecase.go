package category

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type UseCase interface {
	GetCatalogCategories() ([]*models.CategoriesCatalog, error)
	GetSubCategoriesById(categoryId uint64) ([]*models.CategoriesCatalog, error)
}
