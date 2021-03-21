package handler

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools"
	"io/ioutil"
	"net/http"
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

	var productPosition models.ProductPosition
	err = json.Unmarshal(body, &productPosition)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	//err = h.UserUCase.UpdateProfile(currentSession.UserId, &updateUser)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *CartHandler) DeleteProductInCart(w http.ResponseWriter, r *http.Request) {
	h.CartUCase
}

func (h *CartHandler) ChangeProductInCart(w http.ResponseWriter, r *http.Request) {

}

func (h *CartHandler) GetProductsFromCart(w http.ResponseWriter, r *http.Request) {

}

func (h *CartHandler) GetCartDataForOrder(w http.ResponseWriter, r *http.Request) {

}

func (h *CartHandler) AddOrder(w http.ResponseWriter, r *http.Request) {

}
