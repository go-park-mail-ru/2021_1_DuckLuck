package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	session_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/usecase"
	user_usecase "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/usecase"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools"
	"io/ioutil"
	"net/http"
)

type UserHandler struct {
	UserUCase      user_usecase.UserUseCase
	SessionManager session_usecase.UseCase
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.SetJSONResponse(w, "{\"error\": \"can't read body of request\"}", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var authUser models.LoginUser
	err = json.Unmarshal(body, &authUser)
	if err != nil {
		tools.SetJSONResponse(w, "{\"error\": \"can't unmarshal body\"}", http.StatusBadRequest)
		return
	}

	profileUser, err := h.UserUCase.Authorize(&authUser)
	switch err {
	case server_errors.ErrIncorrectUserPassword:
		tools.SetJSONResponse(w, "{\"error\": \"incorrect user password\"}", http.StatusBadRequest)
		return

	case server_errors.ErrIncorrectUserEmail:
		tools.SetJSONResponse(w, "{\"error\": \"incorrect user email\"}", http.StatusBadRequest)
		return
	}

	session, err := h.SessionManager.Create(profileUser.Id)
	if err == server_errors.ErrInternalError {
		tools.SetJSONResponse(w, "{\"error\": \"session wasn't create\"}", http.StatusBadRequest)
		return
	}

	tools.SetCookie(w, models.SessionCookieName, session.Value, models.DurationNewSessionCookie)

	tools.SetJSONResponse(w, "{\"result\": \"success\"}", http.StatusOK)
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.SetJSONResponse(w, "{\"error\": \"can't read body of request\"}", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var newUser models.SignupUser
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		tools.SetJSONResponse(w, "{\"error\": \"can't unmarshal body\"}", http.StatusBadRequest)
		return
	}

	addedUser, err := h.UserUCase.UserRepo.Add(&newUser)
	switch err {
	case server_errors.ErrEmailAlreadyExist:
		tools.SetJSONResponse(w, "{\"error\": \"user email already exist\"}", http.StatusConflict)
		return
	}

	session, err := h.SessionManager.Create(addedUser.Id)
	if err == server_errors.ErrInternalError {
		tools.SetJSONResponse(w, "{\"error\": \"session wasn't create\"}", http.StatusBadRequest)
		return
	}

	tools.SetCookie(w, models.SessionCookieName, session.Value, models.DurationNewSessionCookie)

	tools.SetJSONResponse(w, "{\"result\": \"success\"}", http.StatusCreated)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Middleware auth add session in context
	session, ok := r.Context().Value(models.SessionContextKey).(*models.Session)
	if !ok || session == nil {
		tools.SetJSONResponse(w, "{\"error\": \"session not found in context\"}", http.StatusUnauthorized)
		return
	}

	err := h.SessionManager.DestroyCurrent(session.Value)
	if err == server_errors.ErrInternalError {
		tools.SetJSONResponse(w, "{\"error\": \"session wasn't delete\"}", http.StatusUnauthorized)
		return
	}

	// Auth middleware control existence of session cookie
	sessionCookie, _ := r.Cookie(models.SessionCookieName)
	tools.DestroyCookie(w, sessionCookie)

	tools.SetJSONResponse(w, "{\"result\": \"success\"}", http.StatusOK)
}
