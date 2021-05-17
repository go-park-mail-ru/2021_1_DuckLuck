package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/notification"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/validator"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/logger"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/server_push"
)

type NotificationHandler struct {
	NotificationUCase notification.UseCase
}

func NewHandler(notificationUCase notification.UseCase) notification.Handler {
	return &NotificationHandler{
		NotificationUCase: notificationUCase,
	}
}

func (h *NotificationHandler) SubscribeUser(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("notification_handler", "SubscribeUser", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var credentials models.NotificationCredentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}
	credentials.Sanitize()

	err = validator.ValidateStruct(credentials)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	err = h.NotificationUCase.SubscribeUser(currentSession.UserData.Id, &credentials)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotAddReview, http.StatusBadRequest)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *NotificationHandler) UnsubscribeUser(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("notification_handler", "SubscribeUser", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	err = h.NotificationUCase.UnsubscribeUser(currentSession.UserData.Id)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotAddReview, http.StatusBadRequest)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *NotificationHandler) GetNotificationPublicKey(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("notification_handler", "GetNotificationPublicKey", requireId, err)
		}
	}()

	http_utils.SetJSONResponse(w, models.NotificationPublicKey{Key: server_push.VAPIDPublicKey}, http.StatusOK)
}
