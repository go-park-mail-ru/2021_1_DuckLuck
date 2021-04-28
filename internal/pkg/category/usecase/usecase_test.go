package usecase

import (
	"testing"

	category_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserUseCase_GetSubCategoriesById(t *testing.T) {
	categories := []*models.CategoriesCatalog{
		&models.CategoriesCatalog{
			Id:   4,
			Name: "test",
			Next: nil,
		},
	}
	categoryId := uint64(4)

	t.Run("GetPreviewOrder_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_mock.NewMockRepository(ctrl)
		categoryRepo.
			EXPECT().
			GetNextLevelCategories(categoryId).
			Return(categories, nil)

		userUCase := NewUseCase(categoryRepo)

		userData, err := userUCase.GetSubCategoriesById(categoryId)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, categories, userData, "not equal data")
	})
}
