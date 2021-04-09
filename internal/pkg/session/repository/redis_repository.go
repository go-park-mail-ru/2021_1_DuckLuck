package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	conn *redis.Client
}

func NewSessionRedisRepository(conn *redis.Client) session.Repository {
	return &RedisRepository{
		conn: conn,
	}
}

func (r *RedisRepository) getNewKey(value string) string {
	return fmt.Sprintf("session:%s", value)
}

// Add user session in repository
func (r *RedisRepository) AddSession(session *models.Session) error {
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()
	data, err := json.Marshal(session.UserData)
	if err != nil {
		return errors.ErrCanNotMarshal
	}

	key := r.getNewKey(session.Value)

	fmt.Println(key)
	result, err := r.conn.Set(context.Background(), key, data, models.ExpireSessionCookie).Result()
	if err != nil || result != "OK" {
		return errors.ErrDBInternalError
	}

	return nil
}

// Get user global id from redis db
func (r *RedisRepository) SelectUserIdBySession(sessionValue string) (uint64, error) {
	key := r.getNewKey(sessionValue)

	fmt.Println(key)
	data, err := r.conn.Get(context.Background(), key).Bytes()
	if err != nil {
		return 0, errors.ErrSessionNotFound
	}
	fmt.Println(data)

	userData := &models.UserId{}
	err = json.Unmarshal(data, userData)
	if err != nil {
		return 0, errors.ErrCanNotUnmarshal
	}

	fmt.Println(err)

	return userData.Id, nil
}

// Delete session from db
func (r *RedisRepository) DeleteSessionByValue(sessionValue string) error {
	key := r.getNewKey(sessionValue)

	err := r.conn.Del(context.Background(), key).Err()
	if err != nil {
		return errors.ErrDBInternalError
	}

	return nil
}
