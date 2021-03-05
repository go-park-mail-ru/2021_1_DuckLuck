package middleware

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	session_manager "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/usecase"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools"
)

func Auth(sm *session_manager.UseCase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(models.SessionCookieName)
		if err != nil {
			tools.SetJSONResponse(w, []byte("{\"error\": \"user is unauthorized\"}"), http.StatusUnauthorized)
			return
		}

		sess, err := sm.Check(sessionCookie.Value)
		if err != nil {
			tools.SetJSONResponse(w, []byte("{\"error\": \"user is unauthorized\"}"), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(r.Context(), models.SessionContextKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
