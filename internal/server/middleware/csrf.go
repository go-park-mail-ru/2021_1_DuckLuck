package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/csrf_token"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/logger"
)

func CsrfCheck(u csrf_token.UseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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
				csrfToken := r.Header.Get(models.CsrfTokenHeaderName)
				if csrfToken == "" {
					http_utils.SetJSONResponse(w, errors.ErrNotFoundCsrfToken, http.StatusBadRequest)
					return
				}

				if ok := u.CheckCsrfToken(csrfToken); !ok {
					http_utils.SetJSONResponse(w, errors.ErrIncorrectCsrfToken, http.StatusBadRequest)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
