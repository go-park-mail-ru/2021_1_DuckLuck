package tools

import (
	"io"
	"net/http"
	"os"
	"time"
)

func SetJSONResponse(w http.ResponseWriter, body []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}

func SetFileResponse(w http.ResponseWriter, file *os.File, statusCode int) {
	w.WriteHeader(statusCode)
	io.Copy(w, file)
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
