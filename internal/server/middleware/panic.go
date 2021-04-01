package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools"
)

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requireId := tools.MustGetRequireId(r.Context())
				tools.LogError(r.URL.Path, "middleware", "Panic", requireId, err.(error))
				tools.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
