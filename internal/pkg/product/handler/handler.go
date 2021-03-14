package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	ProductUCase product.UseCase
}

func NewHandler(UCase product.UseCase) product.Handler {
	return &ProductHandler{
		ProductUCase: UCase,
	}
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		tools.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	productById, err := h.ProductUCase.GetProductById(uint64(id))
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
	}

	tools.SetJSONResponse(w, productById, http.StatusOK)
}

func (h *ProductHandler) GetListPreviewProducts(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var paginator models.PaginatorProducts
	err = json.Unmarshal(body, &paginator)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	listPreviewProducts, err := h.ProductUCase.SelectRangeProducts(&paginator)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponse(w, listPreviewProducts, http.StatusOK)
}
