package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	proto "github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/proto/session"

	"google.golang.org/grpc"
)

type SessionUseCase struct {
	SessionClient proto.SessionServiceClient
}

func NewUseCase(sessionConn grpc.ClientConnInterface) session.UseCase {
	return &SessionUseCase{
		SessionClient: proto.NewSessionServiceClient(sessionConn),
	}
}

// Get user id by session value
func (u *SessionUseCase) GetUserIdBySession(sessionCookieValue string) (uint64, error) {
	userId, err := u.SessionClient.GetUserIdBySession(context.Background(), &proto.SessionValue{
		Value: sessionCookieValue,
	})
	if err != nil {
		return 0, errors.ErrSessionNotFound
	}

	return userId.Id, nil
}

// Create new user session and save in repository
func (u *SessionUseCase) CreateNewSession(userId uint64) (*models.Session, error) {
	userSession, err := u.SessionClient.CreateNewSession(context.Background(), &proto.UserId{
		Id: userId,
	})
	if err != nil {
		return nil, errors.ErrInternalError
	}

	return &models.Session{
		Value: userSession.Value.Value,
		UserData: models.UserId{
			Id: userSession.Id.Id,
		},
	}, nil
}

// Destroy session from repository by session value
func (u *SessionUseCase) DestroySession(sessionCookieValue string) error {
	_, err := u.SessionClient.DestroySession(context.Background(), &proto.SessionValue{
		Value: sessionCookieValue,
	})

	if err != nil {
		return errors.ErrSessionNotFound
	}

	return nil
}
