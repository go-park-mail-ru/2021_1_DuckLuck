package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/notification"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	conn *redis.Client
}

func NewSessionRedisRepository(conn *redis.Client) notification.Repository {
	return &RedisRepository{
		conn: conn,
	}
}

func (r *RedisRepository) getNewKey(value uint64) string {
	return fmt.Sprintf("notification:%d", value)
}

func (r *RedisRepository) AddSubscribeUser(userId uint64, subscribes *models.Subscribes) error {
	key := r.getNewKey(userId)

	data, err := json.Marshal(subscribes)
	if err != nil {
		return errors.ErrCanNotMarshal
	}

	err = r.conn.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		return errors.ErrDBInternalError
	}

	return nil
}

func (r *RedisRepository) SelectCredentialsByUserId(userId uint64) (*models.Subscribes, error) {
	subscribes := &models.Subscribes{}
	key := r.getNewKey(userId)

	data, err := r.conn.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, errors.ErrSessionNotFound
	}

	if err = json.Unmarshal(data, subscribes); err != nil {
		return nil, errors.ErrCanNotUnmarshal
	}

	return subscribes, nil
}

func (r *RedisRepository) DeleteSubscribeUser(userId uint64) error {
	key := r.getNewKey(userId)

	err := r.conn.Del(context.Background(), key).Err()
	if err != nil {
		return errors.ErrDBInternalError
	}

	return nil
}
