package middleware

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/logger"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

func Auth(sm session.UseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var err error
			defer func() {
				requireId := http_utils.MustGetRequireId(r.Context())
				if err != nil {
					logger.LogError(r.URL.Path, "middleware", "Auth", requireId, err)
				}
			}()

			sessionCookie, err := r.Cookie(models.SessionCookieName)
			if err != nil {
				http_utils.SetJSONResponse(w, errors.ErrUserUnauthorized, http.StatusUnauthorized)
				return
			}

			sess, err := sm.Check(sessionCookie.Value)
			if err != nil {
				http_utils.SetJSONResponse(w, errors.ErrUserUnauthorized, http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(r.Context(), models.SessionContextKey, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
