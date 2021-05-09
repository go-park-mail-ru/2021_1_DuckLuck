package usecase

import (
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/review/mock"
	user_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReviewUseCase_GetStatisticsByProductId(t *testing.T) {
	productId := uint64(1)
	statistics := models.ReviewStatistics{}

	t.Run("GetStatisticsByProductId_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := mock.NewMockRepository(ctrl)
		reviewRepo.
			EXPECT().
			SelectStatisticsByProductId(productId).
			Return(&statistics, nil)

		userRepo := user_mock.NewMockRepository(ctrl)

		userUCase := NewUseCase(reviewRepo, userRepo)

		result, err := userUCase.GetStatisticsByProductId(productId)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, result, &statistics, "not equal data")
	})
}

func TestReviewUseCase_AddNewReviewForProduct(t *testing.T) {
	userId := uint64(3)
	productId := uint64(3)
	review := models.Review{
		ProductId:     3,
		Rating:        1,
		Advantages:    "asf",
		Disadvantages: "fasf",
		Comment:       "fas",
		IsPublic:      false,
	}

	t.Run("AddNewReviewForProduct_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := mock.NewMockRepository(ctrl)
		reviewRepo.
			EXPECT().
			CheckReview(userId, productId).
			Return(true)
		reviewRepo.
			EXPECT().
			AddReview(userId, &review).
			Return(uint64(1), nil)

		userRepo := user_mock.NewMockRepository(ctrl)

		userUCase := NewUseCase(reviewRepo, userRepo)

		err := userUCase.AddNewReviewForProduct(userId, &review)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("AddNewReviewForProduct_rights_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := mock.NewMockRepository(ctrl)
		reviewRepo.
			EXPECT().
			CheckReview(userId, productId).
			Return(false)

		userRepo := user_mock.NewMockRepository(ctrl)

		userUCase := NewUseCase(reviewRepo, userRepo)

		err := userUCase.AddNewReviewForProduct(userId, &review)

		assert.Equal(t, errors.ErrNoWriteRights, err, "not equal data")
	})

	t.Run("AddNewReviewForProduct_can_no_add_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := mock.NewMockRepository(ctrl)
		reviewRepo.
			EXPECT().
			CheckReview(userId, productId).
			Return(true)
		reviewRepo.
			EXPECT().
			AddReview(userId, &review).
			Return(uint64(1), errors.ErrCanNotAddReview)

		userRepo := user_mock.NewMockRepository(ctrl)

		userUCase := NewUseCase(reviewRepo, userRepo)

		err := userUCase.AddNewReviewForProduct(userId, &review)

		assert.Equal(t, errors.ErrCanNotAddReview, err, "not equal data")
	})
}

func TestReviewUseCase_GetReviewsByProductId(t *testing.T) {
	productId := uint64(1)
	badPaginator := models.PaginatorReviews{
		PageNum:            0,
		Count:              0,
		SortReviewsOptions: models.SortReviewsOptions{},
	}
	reviews := []*models.ViewReview{{}}
	userProfile := models.ProfileUser{}

	t.Run("GetStatisticsByProductId_paginator_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewRepo := mock.NewMockRepository(ctrl)
		userRepo := user_mock.NewMockRepository(ctrl)

		userUCase := NewUseCase(reviewRepo, userRepo)

		_, err := userUCase.GetReviewsByProductId(productId, &badPaginator)
		assert.Equal(t, err, errors.ErrIncorrectPaginator, "not equal data")
	})

	t.Run("GetReviewsByProductId_paginator_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		badPaginator.Count = 2
		badPaginator.PageNum = 4

		reviewRepo := mock.NewMockRepository(ctrl)
		reviewRepo.
			EXPECT().
			GetCountPages(productId, badPaginator.Count).
			Return(1, errors.ErrIncorrectPaginator)

		userRepo := user_mock.NewMockRepository(ctrl)

		userUCase := NewUseCase(reviewRepo, userRepo)

		_, err := userUCase.GetReviewsByProductId(productId, &badPaginator)

		assert.Equal(t, errors.ErrIncorrectPaginator, err, "not equal data")
	})

	t.Run("GetReviewsByProductId_paginator_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		badPaginator.Count = 2
		badPaginator.PageNum = 4

		reviewRepo := mock.NewMockRepository(ctrl)
		reviewRepo.
			EXPECT().
			GetCountPages(productId, badPaginator.Count).
			Return(1, nil)
		reviewRepo.
			EXPECT().
			CreateSortString(badPaginator.SortKey, badPaginator.SortDirection).
			Return("", errors.ErrIncorrectPaginator)

		userRepo := user_mock.NewMockRepository(ctrl)

		userUCase := NewUseCase(reviewRepo, userRepo)

		_, err := userUCase.GetReviewsByProductId(productId, &badPaginator)

		assert.Equal(t, errors.ErrIncorrectPaginator, err, "not equal data")
	})

	t.Run("GetReviewsByProductId_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		badPaginator.Count = 2
		badPaginator.PageNum = 4

		reviewRepo := mock.NewMockRepository(ctrl)

		reviewRepo.
			EXPECT().
			GetCountPages(productId, badPaginator.Count).
			Return(1, nil)

		reviewRepo.
			EXPECT().
			CreateSortString(badPaginator.SortKey, badPaginator.SortDirection).
			Return("", nil)

		reviewRepo.
			EXPECT().
			SelectRangeReviews(productId, "", &badPaginator).
			Return(reviews, nil)

		userRepo := user_mock.NewMockRepository(ctrl)
		userRepo.
			EXPECT().
			SelectProfileById(uint64(0)).
			Return(&userProfile, nil)

		userUCase := NewUseCase(reviewRepo, userRepo)

		_, err := userUCase.GetReviewsByProductId(productId, &badPaginator)

		assert.NoError(t, err, "not error")
	})

	t.Run("GetReviewsByProductId_internal_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		badPaginator.Count = 2
		badPaginator.PageNum = 4

		reviewRepo := mock.NewMockRepository(ctrl)

		reviewRepo.
			EXPECT().
			GetCountPages(productId, badPaginator.Count).
			Return(1, nil)

		reviewRepo.
			EXPECT().
			CreateSortString(badPaginator.SortKey, badPaginator.SortDirection).
			Return("", nil)

		reviewRepo.
			EXPECT().
			SelectRangeReviews(productId, "", &badPaginator).
			Return(reviews, nil)

		userRepo := user_mock.NewMockRepository(ctrl)
		userRepo.
			EXPECT().
			SelectProfileById(uint64(0)).
			Return(&userProfile, errors.ErrInternalError)

		userUCase := NewUseCase(reviewRepo, userRepo)

		_, err := userUCase.GetReviewsByProductId(productId, &badPaginator)

		assert.Error(t, err, "expected error")
	})
}
