package repository

import (
	"fmt"
	"strconv"

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

func (r *RedisRepository) getNewKey(value string) string {
	return fmt.Sprintf("session:%s", value)
}

// Add user session in repository
func (r *RedisRepository) AddSession(session *models.Session) error {
	data := fmt.Sprintf("%d", session.UserId)
	key := r.getNewKey(session.Value)

	result, err := redis.String(r.conn.Do("SET", key, data, "EX", models.ExpireSessionCookie))
	if err != nil || result != "OK" {
		return errors.ErrDBInternalError
	}

	return nil
}

// Get user global id from redis db
func (r *RedisRepository) SelectUserIdBySession(sessionValue string) (uint64, error) {
	key := r.getNewKey(sessionValue)

	data, err := redis.String(r.conn.Do("GET", key))
	if err != nil {
		return 0, errors.ErrSessionNotFound
	}

	userId, err := strconv.ParseUint(data, 10, 64)
	if err != nil {
		return 0, errors.ErrInternalError
	}

	return userId, nil
}

// Delete session from db
func (r *RedisRepository) DeleteSessionByValue(sessionValue string) error {
	key := r.getNewKey(sessionValue)

	result, err := redis.String(r.conn.Do("DEL", key))
	if err != nil || result != "OK" {
		return errors.ErrDBInternalError
	}

	return nil
}
