package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type SessionUseCase struct {
	SessionRepo session.Repository
}

func NewUseCase(SessionRepo session.Repository) session.UseCase {
	return &SessionUseCase{
		SessionRepo: SessionRepo,
	}
}

// Get user id by session value
func (u *SessionUseCase) GetUserIdBySession(sessionCookieValue string) (uint64, error) {
	return u.SessionRepo.SelectUserIdBySession(sessionCookieValue)
}

// Create new user session and save in repository
func (u *SessionUseCase) CreateNewSession(userId uint64) (*models.Session, error) {
	sess := models.NewSession(userId)
	err := u.SessionRepo.AddSession(sess)
	if err != nil {
		return nil, errors.ErrInternalError
	}

	return sess, nil
}

// Destroy session from repository by session value
func (u *SessionUseCase) DestroySession(sessionCookieValue string) error {
	return u.SessionRepo.DeleteSessionByValue(sessionCookieValue)
}
