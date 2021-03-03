package tools

import (
	"net/http"
	"time"
)

func SetJSONResponse(w http.ResponseWriter, jsonStr string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonStr))
	w.WriteHeader(statusCode)
}

func SetCookie(w http.ResponseWriter, cookieValue string, cookieName string, duration time.Duration) {
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
