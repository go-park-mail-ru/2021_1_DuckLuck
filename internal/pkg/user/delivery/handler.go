package delivery

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	session_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/usecase"
	user_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/usecase"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type UserHandler struct {
	UserUCase      user_usecase.UserUseCase
	SessionManager session_usecase.UseCase
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"error\": \"can't read body of request\"}"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var authUser models.LoginUser
	err = json.Unmarshal(body, &authUser)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"error\": \"can't unmarshal body\"}"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profileUser, err := h.UserUCase.Authorize(&authUser)
	switch err {
	case server_errors.ErrIncorrectUserPassword:
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"error\": \"incorrect user password\"}"))
		w.WriteHeader(http.StatusBadRequest)
		return

	case server_errors.ErrIncorrectUserEmail:
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"error\": \"incorrect user email\"}"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := h.SessionManager.Create(profileUser.Id)
	if err == server_errors.ErrInternalError {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"error\": \"session wasn't create\"}"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:    models.SessionCookieName,
		Value:   session.Value,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"result\": \"success\"}"))
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"error\": \"can't read body of request\"}"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var newUser models.SignupUser
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"error\": \"can't unmarshal body\"}"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addedUser, err := h.UserUCase.UserRepo.Add(&newUser)
	switch err {
	case server_errors.ErrEmailAlreadyExist:
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"error\": \"user email already exist\"}"))
		w.WriteHeader(http.StatusConflict)
		return
	}

	session, err := h.SessionManager.Create(addedUser.Id)
	if err == server_errors.ErrInternalError {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"error\": \"session wasn't create\"}"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:    models.SessionCookieName,
		Value:   session.Value,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"result\": \"success\"}"))
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Middleware auth add session in context
	session, ok := r.Context().Value(models.SessionContextKey).(*models.Session)
	if !ok || session == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"error\": \"session not found in context\"}"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := h.SessionManager.DestroyCurrent(session.Value)
	if err == server_errors.ErrInternalError {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"error\": \"session wasn't delete\"}"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookie := http.Cookie{
		Name:    models.SessionCookieName,
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"result\": \"success\"}"))
	w.WriteHeader(http.StatusOK)
}
