package server

import (
	"context"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/pkg/session"
	proto "github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/proto/session"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/server/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

type SessionServer struct {
	SessionUCase session.UseCase
}

func NewSessionServer(sessionUCase session.UseCase) *SessionServer {
	return &SessionServer{
		SessionUCase: sessionUCase,
	}
}

func (s *SessionServer) GetUserIdBySession(ctx context.Context,
	sessionValue *proto.SessionValue) (*proto.UserId, error) {
	userId, err := s.SessionUCase.GetUserIdBySession(sessionValue.Value)
	if err != nil {
		return nil, errors.CreateError(err)
	}

	return &proto.UserId{
		Id: userId,
	}, nil
}

func (s *SessionServer) CreateNewSession(ctx context.Context,
	userId *proto.UserId) (*proto.Session, error) {
	userSession, err := s.SessionUCase.CreateNewSession(userId.Id)
	if err != nil {
		return nil, errors.CreateError(err)
	}

	return &proto.Session{
		Value: &proto.SessionValue{Value: userSession.Value},
		Id:    &proto.UserId{Id: userSession.UserId},
	}, nil
}

func (s *SessionServer) DestroySession(ctx context.Context,
	sessionValue *proto.SessionValue) (*empty.Empty, error) {
	err := s.SessionUCase.DestroySession(sessionValue.Value)
	if err != nil {
		return nil, errors.CreateError(err)
	}

	return &empty.Empty{}, nil
}
