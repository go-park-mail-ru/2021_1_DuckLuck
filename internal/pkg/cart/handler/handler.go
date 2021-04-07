package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/logger"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/validator"
)

type CartHandler struct {
	CartUCase cart.UseCase
}

func NewHandler(cartUCase cart.UseCase) cart.Handler {
	return &CartHandler{
		CartUCase: cartUCase,
	}
}

// Add product in user cart
func (h *CartHandler) AddProductInCart(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError(r.URL.Path, "cart_handler", "AddProductInCart", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	cartArticle := &models.CartArticle{}
	err = json.Unmarshal(body, cartArticle)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = validator.ValidateStruct(cartArticle)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.CartUCase.AddProduct(currentSession.UserId, cartArticle)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// Delete product from user cart
func (h *CartHandler) DeleteProductInCart(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError(r.URL.Path, "cart_handler", "DeleteProductInCart", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	identifier := &models.ProductIdentifier{}
	err = json.Unmarshal(body, identifier)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = validator.ValidateStruct(identifier)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.CartUCase.DeleteProduct(currentSession.UserId, identifier)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// Change product characteristics in user cart
func (h *CartHandler) ChangeProductInCart(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError(r.URL.Path, "cart_handler", "ChangeProductInCart", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	cartArticle := &models.CartArticle{}
	err = json.Unmarshal(body, cartArticle)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = validator.ValidateStruct(cartArticle)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.CartUCase.ChangeProduct(currentSession.UserId, cartArticle)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// Get all preview products from user cart
func (h *CartHandler) GetProductsFromCart(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError(r.URL.Path, "cart_handler", "GetProductsFromCart", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	previewUserCart, err := h.CartUCase.GetPreviewCart(currentSession.UserId)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, previewUserCart, http.StatusOK)
}
