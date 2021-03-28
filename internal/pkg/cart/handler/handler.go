package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools"
)

type CartHandler struct {
	CartUCase cart.UseCase
}

func NewHandler(cartUCase cart.UseCase) cart.Handler {
	return &CartHandler{
		CartUCase: cartUCase,
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

	cartArticle := &models.CartArticle{}
	err = json.Unmarshal(body, cartArticle)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = tools.ValidateStruct(cartArticle)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.CartUCase.AddProduct(currentSession.UserId, cartArticle)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *CartHandler) DeleteProductInCart(w http.ResponseWriter, r *http.Request) {
	currentSession := tools.MustGetSessionFromContext(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	identifier := &models.ProductIdentifier{}
	err = json.Unmarshal(body, identifier)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = tools.ValidateStruct(identifier)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.CartUCase.DeleteProduct(currentSession.UserId, identifier)
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

	cartArticle := &models.CartArticle{}
	err = json.Unmarshal(body, cartArticle)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = tools.ValidateStruct(cartArticle)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.CartUCase.ChangeProduct(currentSession.UserId, cartArticle)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *CartHandler) GetProductsFromCart(w http.ResponseWriter, r *http.Request) {
	currentSession := tools.MustGetSessionFromContext(r.Context())

	previewUserCart, err := h.CartUCase.GetPreviewCart(currentSession.UserId)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponse(w, previewUserCart, http.StatusOK)
}
