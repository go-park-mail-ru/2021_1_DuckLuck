package middleware

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/jwt_token"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/logger"
)

func CsrfCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			requireId := http_utils.MustGetRequireId(r.Context())
			if err != nil {
				logger.LogError(r.URL.Path, "middleware", "CsrfCheck", requireId, err)
			}
		}()

		if r.Method == http.MethodPost || r.Method == http.MethodDelete ||
			r.Method == http.MethodPut || r.Method == http.MethodPatch {
			csrfTokenFromHeader := r.Header.Get(models.CsrfTokenHeaderName)
			if csrfTokenFromHeader == "" {
				http_utils.SetJSONResponse(w, errors.ErrNotFoundCsrfToken, http.StatusBadRequest)
				return
			}

			csrfTokenFromCookie := jwt_token.JwtToken{}
			jwtToken, err := jwt_token.ParseJwtToken(csrfTokenFromHeader, &csrfTokenFromCookie)
			if err != nil || !jwtToken.Valid {
				http_utils.SetJSONResponse(w, errors.ErrIncorrectJwtToken, http.StatusBadRequest)
				return
			}

			t := time.Now()
			if t.After(csrfTokenFromCookie.Expires) {
				http_utils.SetJSONResponse(w, errors.ErrIncorrectJwtToken, http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
