package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type UseCase struct {
	SessionRepo session.Repository
}

func (h *UseCase) Check(sessionCookieValue string) (*models.Session, error) {
	sess, err := h.SessionRepo.GetByValue(sessionCookieValue)
	if err == server_errors.ErrSessionNotFound {
		return nil, server_errors.ErrUserUnauthorized
	}

	return sess, nil
}

func (h *UseCase) Create(userId uint64) (*models.Session, error) {
	sess := models.NewSession(userId)
	err := h.SessionRepo.Add(sess)
	if err != nil {
		return nil, server_errors.ErrInternalError
	}

	return sess, nil
}

func (h *UseCase) DestroyCurrent(sessionCookieValue string) error {
	err := h.SessionRepo.DestroyByValue(sessionCookieValue)
	if err != nil {
		return server_errors.ErrInternalError
	}

	return nil
}
