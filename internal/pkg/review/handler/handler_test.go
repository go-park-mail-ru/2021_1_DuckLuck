package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/review/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReviewHandler_GetReviewStatistics(t *testing.T) {
	productId := uint64(4)
	statistics := models.ReviewStatistics{}

	t.Run("GetReviewStatistics_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewUCase := mock.NewMockUseCase(ctrl)
		reviewUCase.
			EXPECT().
			GetStatisticsByProductId(productId).
			Return(&statistics, nil)

		reviewHandler := NewHandler(reviewUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/review/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		vars := map[string]string{
			"id": fmt.Sprintf("%d", productId),
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.GetReviewStatistics)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})
}

func TestReviewHandler_CheckReviewRights(t *testing.T) {
	productId := uint64(4)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: 3,
		},
	}

	t.Run("CheckReviewRights_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewUCase := mock.NewMockUseCase(ctrl)
		reviewUCase.
			EXPECT().
			CheckReviewUserRights(gomock.Any(), productId).
			Return(nil)

		reviewHandler := NewHandler(reviewUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/review/rights/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		vars := map[string]string{
			"id": fmt.Sprintf("%d", productId),
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.CheckReviewRights)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("AddNewReview_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUCase := mock.NewMockUseCase(ctrl)

		productHandler := NewHandler(productUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/review/rights/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(productHandler.CheckReviewRights)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

}

func TestReviewHandler_AddNewReview(t *testing.T) {
	//productId := uint64(4)
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}
	review := models.Review{
		ProductId:     4,
		Rating:        0,
		Advantages:    "",
		Disadvantages: "",
		Comment:       "",
		IsPublic:      false,
	}

	t.Run("AddNewReview_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reviewUCase := mock.NewMockUseCase(ctrl)
		reviewUCase.
			EXPECT().
			AddNewReviewForProduct(userId, &review).
			Return(nil)

		reviewHandler := NewHandler(reviewUCase)

		reviewBytes, _ := json.Marshal(review)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/review/product",
			bytes.NewBuffer(reviewBytes))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.AddNewReview)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("AddNewReview_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUCase := mock.NewMockUseCase(ctrl)

		productHandler := NewHandler(productUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/review/rights/product",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(productHandler.CheckReviewRights)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

}
