package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/logger"
)

type OrderHandler struct {
	OrderUCase order.UseCase
	CartUCase  cart.UseCase
}

func NewHandler(orderUCase order.UseCase, cartUCase cart.UseCase) order.Handler {
	return &OrderHandler{
		OrderUCase: orderUCase,
		CartUCase:  cartUCase,
	}
}

func (h *OrderHandler) GetOrderFromCart(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError(r.URL.Path, "order_handler", "GetOrderFromCart", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	previewCart, err := h.CartUCase.GetPreviewCart(currentSession.UserData.Id)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	previewOrder, err := h.OrderUCase.GetPreviewOrder(currentSession.UserData.Id, previewCart)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, previewOrder, http.StatusOK)
}

func (h *OrderHandler) AddCompletedOrder(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError(r.URL.Path, "order_handler", "AddCompletedOrder", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var userOrder models.Order
	err = json.Unmarshal(body, &userOrder)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}
	userOrder.Sanitize()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	previewCart, err := h.CartUCase.GetPreviewCart(currentSession.UserData.Id)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	_, err = h.OrderUCase.AddCompletedOrder(&userOrder, currentSession.UserData.Id, previewCart)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError(r.URL.Path, "order_handler", "GetUserOrders", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	orders, err := h.OrderUCase.GetOrders(currentSession.UserData.Id)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, orders, http.StatusOK)
}
