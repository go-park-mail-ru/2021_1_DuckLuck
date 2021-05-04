package handler

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/logger"

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

// Get preview category catalog
func (h *CategoryHandler) GetCatalogCategories(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("category_handler", "GetCatalogCategories", requireId, err)
		}
	}()

	categories, err := h.CategoryUCase.GetCatalogCategories()
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, categories, http.StatusOK)
}

// Get subcategories of category
func (h *CategoryHandler) GetSubCategories(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("category_handler", "GetSubCategories", requireId, err)
		}
	}()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	categories, err := h.CategoryUCase.GetSubCategoriesById(uint64(id))
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, categories, http.StatusOK)
}
