package middleware

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/logger"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requireId := http_utils.MustGetRequireId(r.Context())
				logger.LogError(r.URL.Path, "middleware", "Panic", requireId, err.(error))
				http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
