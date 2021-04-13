package category

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category UseCase

type UseCase interface {
	GetCatalogCategories() ([]*models.CategoriesCatalog, error)
	GetSubCategoriesById(categoryId uint64) ([]*models.CategoriesCatalog, error)
}
