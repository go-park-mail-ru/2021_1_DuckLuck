package handler

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/csrf_token"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
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
	jwtToken, err := jwt_token.CreateJwtToken([]byte(csrfToken.Value),
		time.Now().Add(models.ExpireCsrfToken*time.Second))

	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}
	csrfToken.Value = jwtToken

	http_utils.SetJSONResponse(w, csrfToken, http.StatusOK)
}
