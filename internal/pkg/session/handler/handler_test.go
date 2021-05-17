package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/mock"

	"github.com/golang/mock/gomock"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
)

func TestSessionHandler_CheckSession(t *testing.T) {
	t.Run("CheckSession_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionUCase := mock.NewMockUseCase(ctrl)
		sessionHandler := NewHandler(sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/session",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sessionHandler.CheckSession)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})
}
