package handler

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/logger"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/validator"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
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
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError(r.URL.Path, "product_handler", "GetProduct", requireId, err)
		}
	}()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	productById, err := h.ProductUCase.GetProductById(uint64(id))
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
	}

	http_utils.SetJSONResponse(w, productById, http.StatusOK)
}

func (h *ProductHandler) GetListPreviewProducts(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError(r.URL.Path, "product_handler", "GetListPreviewProducts", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var paginator models.PaginatorProducts
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

	listPreviewProducts, err := h.ProductUCase.SelectRangeProducts(&paginator)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, listPreviewProducts, http.StatusOK)
}
