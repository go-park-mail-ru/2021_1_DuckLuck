package repository

import (
	"encoding/json"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/gomodule/redigo/redis"
)

type RedisRepository struct {
	conn redis.Conn
}

func NewSessionRedisRepository(conn redis.Conn) session.Repository {
	return &RedisRepository{
		conn: conn,
	}
}

func (rr *RedisRepository) AddSession(session *models.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return errors.ErrCanNotMarshal
	}

	_, err = redis.String(rr.conn.Do("SET", session.Value, data, "EX", models.ExpireSessionCookie))
	if err != nil {
		return errors.ErrDBInternalError
	}
	return nil
}

func (rr *RedisRepository) SelectSessionByValue(sessionValue string) (*models.Session, error) {
	data, err := redis.Bytes(rr.conn.Do("GET", sessionValue))
	if err != nil {
		return nil, errors.ErrSessionNotFound
	}

	currentSession := &models.Session{}
	err = json.Unmarshal(data, currentSession)
	if err != nil {
		return nil, errors.ErrCanNotUnmarshal
	}

	return currentSession, nil
}

func (rr *RedisRepository) DeleteByValue(sessionValue string) error {
	_, err := redis.Int(rr.conn.Do("DEL", sessionValue))
	if err != nil {
		return errors.ErrDBInternalError
	}
	return nil
}
