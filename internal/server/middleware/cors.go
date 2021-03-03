package middleware

import (
	"net/http"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, "+
			"Accept-Encoding, Authorization")

		// If method of request is only for get options
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}
