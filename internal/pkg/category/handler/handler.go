package handler

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools"

	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	CategoryUCase category.UseCase
}

func NewHandler(UCase category.UseCase) category.Handler {
	return &CategoryHandler{
		CategoryUCase: UCase,
	}
}

func (h *CategoryHandler) GetCatalogCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.CategoryUCase.GetCatalogCategories()
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponse(w, categories, http.StatusOK)
}

func (h *CategoryHandler) GetSubCategories(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		tools.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	categories, err := h.CategoryUCase.GetSubCategoriesById(uint64(id))
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponse(w, categories, http.StatusOK)
}
