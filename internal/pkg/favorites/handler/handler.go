package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/favorites"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/validator"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/logger"

	"github.com/gorilla/mux"
)

type FavoritesHandler struct {
	FavoritesUCase favorites.UseCase
}

func NewHandler(favoritesUCase favorites.UseCase) favorites.Handler {
	return &FavoritesHandler{
		FavoritesUCase: favoritesUCase,
	}
}

func (h *FavoritesHandler) AddProductToFavorites(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("favorites_handler", "AddProductToFavorites", requireId, err)
		}
	}()

	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])
	if err != nil || productId < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	if err = h.FavoritesUCase.AddProductToFavorites(uint64(productId), currentSession.UserData.Id); err != nil {
		http_utils.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *FavoritesHandler) DeleteProductFromFavorites(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("favorites_handler", "DeleteProductFromFavorites", requireId, err)
		}
	}()

	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])
	if err != nil || productId < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	if err = h.FavoritesUCase.DeleteProductFromFavorites(uint64(productId), currentSession.UserData.Id); err != nil {
		http_utils.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *FavoritesHandler) GetListPreviewFavorites(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("favorites_handler", "GetListPreviewFavorites", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var paginator models.PaginatorFavorites
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

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	listPreviewFavorites, err := h.FavoritesUCase.GetRangeFavorites(&paginator, currentSession.UserData.Id)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, listPreviewFavorites, http.StatusOK)
}

func (h *FavoritesHandler) GetUserFavorites(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("favorites_handler", "GetUserFavorites", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	listFavorites, err := h.FavoritesUCase.GetUserFavorites(currentSession.UserData.Id)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, listFavorites, http.StatusOK)
}
