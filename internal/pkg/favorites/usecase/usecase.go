package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/favorites"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type FavoritesUseCase struct {
	FavoritesRepo favorites.Repository
}

func NewUseCase(favoritesRepo favorites.Repository) favorites.UseCase {
	return &FavoritesUseCase{
		FavoritesRepo: favoritesRepo,
	}
}

func (u *FavoritesUseCase) AddProductToFavorites(productId, userId uint64) error {
	return u.FavoritesRepo.AddProductToFavorites(productId, userId)
}

func (u *FavoritesUseCase) DeleteProductFromFavorites(productId, userId uint64) error {
	return u.FavoritesRepo.DeleteProductFromFavorites(productId, userId)
}

func (u *FavoritesUseCase) GetRangeFavorites(paginator *models.PaginatorFavorites,
	userId uint64) (*models.RangeFavorites, error) {
	if paginator.PageNum < 1 || paginator.Count < 1 {
		return nil, errors.ErrIncorrectPaginator
	}

	// Max count pages
	countPages, err := u.FavoritesRepo.GetCountPages(userId, paginator.Count)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Keys for sort items
	sortString, err := u.FavoritesRepo.CreateSortString(paginator.SortKey, paginator.SortDirection)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Get range of favorites
	products, err := u.FavoritesRepo.SelectRangeFavorites(paginator, sortString, userId)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	return &models.RangeFavorites{
		ListPreviewProducts: products,
		MaxCountPages:       countPages,
	}, nil
}

func (u *FavoritesUseCase) GetUserFavorites(userId uint64) (*models.UserFavorites, error) {
	return u.FavoritesRepo.GetUserFavorites(userId)
}
