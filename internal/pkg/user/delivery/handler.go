package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"io/ioutil"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	session_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/usecase"
	user_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/usecase"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools"
)

type UserHandler struct {
	UserUCase      user_usecase.UserUseCase
	SessionManager session_usecase.UseCase
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"can't read body of request\"}"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var authUser models.LoginUser
	err = json.Unmarshal(body, &authUser)
	if err != nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"can't unmarshal body\"}"), http.StatusBadRequest)
		return
	}

	profileUser, err := h.UserUCase.Authorize(&authUser)
	switch err {
	case server_errors.ErrIncorrectUserPassword:
		tools.SetJSONResponse(w, []byte("{\"error\": \"incorrect user password\"}"), http.StatusBadRequest)
		return

	case server_errors.ErrIncorrectUserEmail:
		tools.SetJSONResponse(w, []byte("{\"error\": \"incorrect user email\"}"), http.StatusBadRequest)
		return
	}

	session, err := h.SessionManager.Create(profileUser.Id)
	if err == server_errors.ErrInternalError {
		tools.SetJSONResponse(w, []byte("{\"error\": \"session wasn't create\"}"), http.StatusBadRequest)
		return
	}

	tools.SetCookie(w, models.SessionCookieName, session.Value, models.DurationNewSessionCookie)

	tools.SetJSONResponse(w, []byte("{\"result\": \"success\"}"), http.StatusOK)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"can't read body of request\"}"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var profileUser models.ProfileUser
	err = json.Unmarshal(body, &profileUser)
	if err != nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"can't unmarshal body\"}"), http.StatusBadRequest)
		return
	}

	err = h.UserUCase.UserRepo.Update(&profileUser)
	if err != nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"can't update user\"}"), http.StatusBadRequest)
		return
	}

	tools.SetJSONResponse(w, []byte("{\"result\": \"success\"}"), http.StatusOK)
}

func (h *UserHandler) UpdateProfileAvatar(w http.ResponseWriter, r *http.Request) {
	// Middleware auth add session in context
	session, ok := r.Context().Value(models.SessionContextKey).(*models.Session)
	if !ok || session == nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"session not found in context\"}"), http.StatusUnauthorized)
		return
	}

	fileName, err := tools.UploadFile(r, "avatar")
	switch err {
	case errors.ErrServerSystem:
		tools.SetJSONResponse(w, []byte("{\"error\": \"system error\"}"), http.StatusInternalServerError)
		return
	case errors.ErrFileNotRead:
		tools.SetJSONResponse(w, []byte("{\"error\": \"file can't read\"}"), http.StatusInternalServerError)
		return
	}

	fileUrl, err := h.UserUCase.SetAvatar(session.UserId, fileName)
	switch err {
	case errors.ErrServerSystem:
		tools.SetJSONResponse(w, []byte("{\"error\": \"system error\"}"), http.StatusInternalServerError)
		return
	case server_errors.ErrUserNotFound:
		tools.SetJSONResponse(w, []byte("{\"error\": \"user not found\"}"), http.StatusBadRequest)
		return
	}

	tools.SetJSONResponse(w, []byte(fmt.Sprintf("{\"result\": \"%s\"}", fileUrl)), http.StatusOK)
}

func (h *UserHandler) GetProfileAvatar(w http.ResponseWriter, r *http.Request) {
	// Middleware auth add session in context
	session, ok := r.Context().Value(models.SessionContextKey).(*models.Session)
	if !ok || session == nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"session not found in context\"}"), http.StatusUnauthorized)
		return
	}

	fileName, err := h.UserUCase.GetAvatar(session.UserId)
	if err == server_errors.ErrUserNotFound {
		tools.SetJSONResponse(w, []byte("{\"error\": \"user not found\"}"), http.StatusInternalServerError)
		return
	}

	tools.SetJSONResponse(w, []byte(fmt.Sprintf("{\"result\": \"%s\"}", fileName)), http.StatusOK)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(models.SessionContextKey).(*models.Session)
	if !ok || session == nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"session not found in context\"}"), http.StatusUnauthorized)
		return
	}

	profileUser, err := h.UserUCase.UserRepo.GetById(session.UserId)
	if err == server_errors.ErrUserNotFound {
		tools.SetJSONResponse(w, []byte("{\"error\": \"user not found\"}"), http.StatusBadRequest)
		return
	}

	// Url to avatar
	profileUser.Avatar = configs.UrlToAvatar + profileUser.Avatar

	result, err := json.Marshal(profileUser)
	if err != nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"can't marshal body\"}"), http.StatusBadRequest)
		return
	}

	tools.SetJSONResponse(w, result, http.StatusOK)
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"can't read body of request\"}"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var newUser models.SignupUser
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"can't unmarshal body\"}"), http.StatusBadRequest)
		return
	}

	addedUser, err := h.UserUCase.UserRepo.Add(&newUser)
	switch err {
	case server_errors.ErrEmailAlreadyExist:
		tools.SetJSONResponse(w, []byte("{\"error\": \"user email already exist\"}"), http.StatusConflict)
		return
	}

	session, err := h.SessionManager.Create(addedUser.Id)
	if err == server_errors.ErrInternalError {
		tools.SetJSONResponse(w, []byte("{\"error\": \"session wasn't create\"}"), http.StatusBadRequest)
		return
	}

	tools.SetCookie(w, models.SessionCookieName, session.Value, models.DurationNewSessionCookie)

	tools.SetJSONResponse(w, []byte("{\"result\": \"success\"}"), http.StatusCreated)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Middleware auth add session in context
	session, ok := r.Context().Value(models.SessionContextKey).(*models.Session)
	if !ok || session == nil {
		tools.SetJSONResponse(w, []byte("{\"error\": \"session not found in context\"}"), http.StatusUnauthorized)
		return
	}

	err := h.SessionManager.DestroyCurrent(session.Value)
	if err == server_errors.ErrInternalError {
		tools.SetJSONResponse(w, []byte("{\"error\": \"session wasn't delete\"}"), http.StatusUnauthorized)
		return
	}

	// Auth middleware control existence of session cookie
	sessionCookie, _ := r.Cookie(models.SessionCookieName)
	tools.DestroyCookie(w, sessionCookie)

	tools.SetJSONResponse(w, []byte("{\"result\": \"success\"}"), http.StatusOK)
}
