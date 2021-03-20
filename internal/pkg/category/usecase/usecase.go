package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type CategoryUseCase struct {
	CategoryRepo category.Repository
}

func NewUseCase(repo category.Repository) category.UseCase {
	return &CategoryUseCase{
		CategoryRepo: repo,
	}
}

func (lr *CategoryUseCase) GetCatalogCategories() ([]*models.CategoriesCatalog, error) {
	categories, err := lr.CategoryRepo.GetCategoriesByLevel(1)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}

	for _, category := range categories {
		nextLevel, err := lr.CategoryRepo.GetNextLevelCategories(category.Id)
		if err != nil {
			return nil, errors.ErrDBInternalError
		}
		category.Next = nextLevel

		for _, subCategory := range category.Next {
			nextLevel, err = lr.CategoryRepo.GetNextLevelCategories(subCategory.Id)
			if err != nil {
				return nil, errors.ErrDBInternalError
			}
			subCategory.Next = nextLevel
		}
	}

	return categories, nil
}

func (lr *CategoryUseCase) GetSubCategoriesById(categoryId uint64) ([]*models.CategoriesCatalog, error) {
	categories, err := lr.CategoryRepo.GetNextLevelCategories(categoryId)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}

	return categories, nil
}
