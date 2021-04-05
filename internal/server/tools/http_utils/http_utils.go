package http_utils

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

func SetJSONResponse(w http.ResponseWriter, body interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	result, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"can't marshal body\"}"))
		return
	}
	w.WriteHeader(statusCode)
	w.Write(result)
}

func SetJSONResponseSuccess(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte("{\"error\": \"success\"}"))
}

func SetCookie(w http.ResponseWriter, cookieName string, cookieValue string, duration time.Duration) {
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		Expires:  time.Now().Add(duration),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func DestroyCookie(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	cookie.Path = "/"
	http.SetCookie(w, cookie)
}

func MustGetSessionFromContext(ctx context.Context) *models.Session {
	session, ok := ctx.Value(models.SessionContextKey).(*models.Session)
	if !ok || session == nil {
		panic(errors.ErrSessionNotFound.Error())
	}

	return session
}

func MustGetRequireId(ctx context.Context) string {
	requireId, ok := ctx.Value(models.RequireIdKey).(string)
	if !ok {
		panic(errors.ErrRequireIdNotFound.Error())
	}

	return requireId
}
