package handler

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/csrf_token"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/hasher"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/jwt_token"
)

type CsrfTokenHandler struct{}

func NewHandler() csrf_token.Handler {
	return &CsrfTokenHandler{}
}

// Get new csrf token for client
func (h *CsrfTokenHandler) GetCsrfToken(w http.ResponseWriter, r *http.Request) {
	csrfToken := models.NewCsrfToken()
	hash, err := hasher.GenerateHashFromCsrfToken(csrfToken.Value)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	jwtToken, err := jwt_token.CreateJwtToken(hash, time.Now().Add(models.ExpireCsrfToken*time.Second))
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetCookie(w, models.CsrfTokenCookieName, jwtToken, models.ExpireCsrfToken*time.Second)
	http_utils.SetJSONResponse(w, csrfToken, http.StatusOK)
}
