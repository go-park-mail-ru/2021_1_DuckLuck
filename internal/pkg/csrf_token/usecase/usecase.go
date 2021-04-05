package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/csrf_token"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type CsrfTokenUseCase struct {
	CsrfTokenRepo csrf_token.Repository
}

func NewUseCase(csrfTokenRepo csrf_token.Repository) csrf_token.UseCase {
	return &CsrfTokenUseCase{
		CsrfTokenRepo: csrfTokenRepo,
	}
}

// Create new csrf token
func (u *CsrfTokenUseCase) CreateCsrfToken() (*models.CsrfToken, error) {
	token := models.NewCsrfToken()
	err := u.CsrfTokenRepo.AddCsrfToken(token.Value)
	if err != nil {
		return nil, errors.ErrInternalError
	}

	return token, nil
}

// Check csrf token
func (u *CsrfTokenUseCase) CheckCsrfToken(csrfTokenValue string) bool {
	return u.CsrfTokenRepo.CheckCsrfToken(csrfTokenValue)
}
