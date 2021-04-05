package handler

import (
	"net/http"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/csrf_token"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/http_utils"
)

type CsrfTokenHandler struct {
	CsrfTokenUCase csrf_token.UseCase
}

func NewHandler(UCase csrf_token.UseCase) csrf_token.Handler {
	return &CsrfTokenHandler{
		CsrfTokenUCase: UCase,
	}
}

// Get new csrf token for client
func (h *CsrfTokenHandler) GetCsrfToken(w http.ResponseWriter, r *http.Request) {
	csrfToken, err := h.CsrfTokenUCase.CreateCsrfToken()
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, csrfToken, http.StatusOK)
}
