package delivery

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	ProductUCase product.UseCase
	ProductRepo  product.Repository
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		tools.SetJSONResponse(w, []byte("{\"error\": \"incorrect product id\"}"), http.StatusInternalServerError)
		return
	}

	product, err := h.ProductRepo.GetById(uint64(id))
	if err == server_errors.ErrProductNotFound {
		tools.SetJSONResponse(w, []byte("{\"error\": \"product not found\"}"), http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(product)
	if err != nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"can't marshal body\"}"), http.StatusBadRequest)
		return
	}

	tools.SetJSONResponse(w, result, http.StatusOK)
}

func (h *ProductHandler) GetListPreviewProducts(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"can't read body of request\"}"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var paginator models.PaginatorProducts
	err = json.Unmarshal(body, &paginator)
	if err != nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"can't unmarshal body\"}"), http.StatusBadRequest)
		return
	}

	listPreviewProducts, err := h.ProductRepo.GetListPreviewProducts(&paginator)
	if err == server_errors.ErrUserNotFound {
		tools.SetJSONResponse(w, []byte("{\"error\": \"user not found\"}"), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(listPreviewProducts)
	if err != nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"can't marshal body\"}"), http.StatusBadRequest)
		return
	}

	tools.SetJSONResponse(w, result, http.StatusOK)
}
