package handler

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools"
	"io/ioutil"
	"net/http"
	"strconv"
)

type CartHandler struct {
	CartUCase cart.UseCase
}

func NewHandler(UCase cart.UseCase) cart.Handler {
	return &CartHandler{
		CartUCase: UCase,
	}
}

func (h *CartHandler) AddProductInCart(w http.ResponseWriter, r *http.Request) {
	currentSession := tools.MustGetSessionFromContext(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	productPosition := &models.ProductPosition{}
	err = json.Unmarshal(body, productPosition)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = h.CartUCase.AddProduct(currentSession.UserId, productPosition)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *CartHandler) DeleteProductInCart(w http.ResponseWriter, r *http.Request) {
	currentSession := tools.MustGetSessionFromContext(r.Context())

	err := r.ParseForm()
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	str := r.Form["id"][0]
	productId, err := strconv.Atoi(str)
	if err != nil || productId < 0 {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.CartUCase.DeleteProduct(currentSession.UserId, uint64(productId))
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *CartHandler) ChangeProductInCart(w http.ResponseWriter, r *http.Request) {
	currentSession := tools.MustGetSessionFromContext(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	productPosition := &models.ProductPosition{}
	err = json.Unmarshal(body, productPosition)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = h.CartUCase.ChangeProduct(currentSession.UserId, productPosition)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *CartHandler) GetProductsFromCart(w http.ResponseWriter, r *http.Request) {
	currentSession := tools.MustGetSessionFromContext(r.Context())

	userCart, err := h.CartUCase.GetCart(currentSession.UserId)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponse(w, *userCart, http.StatusOK)
}
