package repository

import (
	"fmt"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/csrf_token"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/gomodule/redigo/redis"
)

type RedisRepository struct {
	conn redis.Conn
}

func NewCsrfTokenRedisRepository(conn redis.Conn) csrf_token.Repository {
	return &RedisRepository{
		conn: conn,
	}
}

func (r *RedisRepository) getNewKey(value string) string {
	return fmt.Sprintf("csrf_token:%s", value)
}

// Add csrf token in db
func (r *RedisRepository) AddCsrfToken(tokenValue string) error {
	key := r.getNewKey(tokenValue)

	result, err := redis.String(r.conn.Do("SET", key, 1, "EX", models.ExpireCsrfToken))
	if err != nil || result != "OK" {
		return errors.ErrDBInternalError
	}

	return nil
}

// Check csrf token in db
func (r *RedisRepository) CheckCsrfToken(tokenValue string) bool {
	key := r.getNewKey(tokenValue)

	result, err := redis.String(r.conn.Do("GET", key))
	if err != nil || result == "" {
		return false
	}

	return true
}
