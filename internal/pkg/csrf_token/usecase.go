package csrf_token

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type UseCase interface {
	CreateCsrfToken() (*models.CsrfToken, error)
	CheckCsrfToken(csrfTokenValue string) bool
}
