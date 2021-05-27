package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/admin"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/validator"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/logger"
)

type AdminHandler struct {
	AdminUCase admin.UseCase
}

func NewHandler(adminUCase admin.UseCase) admin.Handler {
	return &AdminHandler{
		AdminUCase: adminUCase,
	}
}

func (h *AdminHandler) ChangeOrderStatus(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("admin_handler", "ChangeOrderStatus", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var updateOrder models.UpdateOrder
	err = json.Unmarshal(body, &updateOrder)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}
	updateOrder.Sanitize()

	err = validator.ValidateStruct(updateOrder)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.AdminUCase.ChangeOrderStatus(&updateOrder)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}
