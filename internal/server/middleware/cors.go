package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", configs.CorsOrigins)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// If method of request is only for get options
		if r.Method == http.MethodOptions {
			w.Header().Add("Content-Type", "text/plain")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, "+
				"Accept-Encoding, Authorization")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
