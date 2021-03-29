package handler

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools"
)

type UserHandler struct {
	UserUCase    user.UseCase
	SessionUCase session.UseCase
}

func NewHandler(userUCase user.UseCase, sessionUCase session.UseCase) user.Handler {
	return &UserHandler{
		UserUCase:    userUCase,
		SessionUCase: sessionUCase,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var authUser models.LoginUser
	err = json.Unmarshal(body, &authUser)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = tools.ValidateStruct(authUser)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	profileUser, err := h.UserUCase.Authorize(&authUser)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	currentSession, err := h.SessionUCase.Create(profileUser.Id)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetCookie(w, models.SessionCookieName, currentSession.Value, models.ExpireSessionCookie*time.Second)
	tools.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	currentSession := tools.MustGetSessionFromContext(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var updateUser models.UpdateUser
	err = json.Unmarshal(body, &updateUser)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = tools.ValidateStruct(updateUser)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.UserUCase.UpdateProfile(currentSession.UserId, &updateUser)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponseSuccess(w, http.StatusOK)
}

func (h *UserHandler) UpdateProfileAvatar(w http.ResponseWriter, r *http.Request) {
	currentSession := tools.MustGetSessionFromContext(r.Context())

	fileName, err := tools.UploadFile(r, "avatar", configs.PathToUpload+configs.UrlToAvatar)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	fileUrl, err := h.UserUCase.SetAvatar(currentSession.UserId, fileName)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponse(w, models.Avatar{Url: sql.NullString{String: fileUrl}}, http.StatusOK)
}

func (h *UserHandler) GetProfileAvatar(w http.ResponseWriter, r *http.Request) {
	currentSession := tools.MustGetSessionFromContext(r.Context())

	fileUrl, err := h.UserUCase.GetAvatar(currentSession.UserId)
	if err == errors.ErrUserNotFound {
		tools.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponse(w, models.Avatar{Url: sql.NullString{String: fileUrl}}, http.StatusOK)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	currentSession := tools.MustGetSessionFromContext(r.Context())

	profileUser, err := h.UserUCase.GetUserById(currentSession.UserId)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	tools.SetJSONResponse(w, profileUser, http.StatusOK)
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var newUser models.SignupUser
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		tools.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = tools.ValidateStruct(newUser)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	addedUserId, err := h.UserUCase.AddUser(&newUser)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusConflict)
		return
	}

	currentSession, err := h.SessionUCase.Create(addedUserId)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	tools.SetCookie(w, models.SessionCookieName, currentSession.Value, models.ExpireSessionCookie*time.Second)
	tools.SetJSONResponseSuccess(w, http.StatusCreated)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Middleware auth add session in context
	currentSession := tools.MustGetSessionFromContext(r.Context())

	err := h.SessionUCase.DestroyCurrent(currentSession.Value)
	if err != nil {
		tools.SetJSONResponse(w, errors.CreateError(err), http.StatusUnauthorized)
		return
	}

	// Auth middleware control existence of session cookie
	sessionCookie, _ := r.Cookie(models.SessionCookieName)
	tools.DestroyCookie(w, sessionCookie)

	tools.SetJSONResponseSuccess(w, http.StatusOK)
}
