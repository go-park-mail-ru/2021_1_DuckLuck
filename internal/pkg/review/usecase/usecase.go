package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/review"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type ReviewUseCase struct {
	ReviewRepo review.Repository
	UserRepo   user.Repository
}

func NewUseCase(reviewRepo review.Repository, userRepo user.Repository) review.UseCase {
	return &ReviewUseCase{
		ReviewRepo: reviewRepo,
		UserRepo:   userRepo,
	}
}

func (u *ReviewUseCase) GetStatisticsByProductId(productId uint64) (*models.ReviewStatistics, error) {
	return u.ReviewRepo.SelectStatisticsByProductId(productId)
}

func (u *ReviewUseCase) CheckReviewUserRights(userId uint64, productId uint64) error {
	rights := u.ReviewRepo.CheckReview(userId, productId)
	if !rights {
		return errors.ErrNoWriteRights
	}

	return nil
}

func (u *ReviewUseCase) AddNewReviewForProduct(userId uint64, review *models.Review) error {
	rights := u.ReviewRepo.CheckReview(userId, uint64(review.ProductId))
	if !rights {
		return errors.ErrNoWriteRights
	}

	_, err := u.ReviewRepo.AddReview(userId, review)
	if err != nil {
		return errors.ErrCanNotAddReview
	}
	return nil
}

func (u *ReviewUseCase) GetReviewsByProductId(productId uint64,
	paginator *models.PaginatorReviews) (*models.RangeReviews, error) {
	if paginator.PageNum < 1 || paginator.Count < 1 {
		return nil, errors.ErrIncorrectPaginator
	}

	// Max count pages in catalog
	countPages, err := u.ReviewRepo.GetCountPages(productId, paginator.Count)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Keys for sort reviews in catalog
	sortString, err := u.ReviewRepo.CreateSortString(paginator.SortKey, paginator.SortDirection)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Get range of reviews
	reviews, err := u.ReviewRepo.SelectRangeReviews(productId, sortString, paginator)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Get user data for review
	for _, userReview := range reviews {
		userInfo, err := u.UserRepo.SelectProfileById(uint64(userReview.UserId))
		if err != nil {
			return nil, errors.ErrInternalError
		}

		if userReview.IsPublic {
			userReview.UserAvatar = userInfo.Avatar.Url
			userReview.UserName = fmt.Sprintf("%s %s", userInfo.FirstName, userInfo.LastName)
		}
	}

	return &models.RangeReviews{
		ListPreviews:  reviews,
		MaxCountPages: countPages,
	}, nil
}
