package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/review"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/validator"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/logger"

	"github.com/gorilla/mux"
)

type ReviewHandler struct {
	ReviewUCase review.UseCase
}

func NewHandler(reviewUCase review.UseCase) review.Handler {
	return &ReviewHandler{
		ReviewUCase: reviewUCase,
	}
}

// Get statistics about reviews for product
func (h *ReviewHandler) GetReviewStatistics(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("review_handler", "GetReviewInfo", requireId, err)
		}
	}()

	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])
	if err != nil || productId < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	productById, err := h.ReviewUCase.GetStatisticsByProductId(uint64(productId))
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, productById, http.StatusOK)
}

// Check rights for write new review
func (h *ReviewHandler) CheckReviewRights(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("review_handler", "AddCompletedOrder", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])
	if err != nil || productId < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	err = h.ReviewUCase.CheckReviewUserRights(currentSession.UserData.Id, uint64(productId))
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// Add new review for product
func (h *ReviewHandler) AddNewReview(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("review_handler", "AddNewReview", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var userReview models.Review
	err = json.Unmarshal(body, &userReview)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}
	userReview.Sanitize()

	err = validator.ValidateStruct(userReview)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	err = h.ReviewUCase.AddNewReviewForProduct(currentSession.UserData.Id, &userReview)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotAddReview, http.StatusBadRequest)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// Get all reviews for product
func (h *ReviewHandler) GetReviewsForProduct(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("review_handler", "GetReviewsForProduct", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var paginator models.PaginatorReviews
	err = json.Unmarshal(body, &paginator)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = validator.ValidateStruct(paginator)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])
	if err != nil || productId < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	reviews, err := h.ReviewUCase.GetReviewsByProductId(uint64(productId), &paginator)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, reviews, http.StatusOK)
}
