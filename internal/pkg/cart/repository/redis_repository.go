package repository

import (
	"encoding/json"
	"fmt"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/gomodule/redigo/redis"
)

type RedisRepository struct {
	conn redis.Conn
}

func NewSessionRedisRepository(conn redis.Conn) cart.Repository {
	return &RedisRepository{
		conn: conn,
	}
}

func (rr *RedisRepository) GetNewKey(value uint64) string {
	return fmt.Sprintf("cart:%d", value)
}

func (rr *RedisRepository) DeleteCart(userId uint64) error {
	key := rr.GetNewKey(userId)

	_, err := redis.String(rr.conn.Do("DEL", key))
	if err != nil {
		return errors.ErrDBInternalError
	}
	return nil
}

func (rr *RedisRepository) GetCart(userId uint64) (*models.Cart, error) {
	userCart := &models.Cart{}
	key := rr.GetNewKey(userId)

	data, err := redis.Bytes(rr.conn.Do("GET", key))
	if err != nil {
		return nil, errors.ErrCartNotFound
	}

	if err = json.Unmarshal(data, &userCart); err != nil {
		return nil, errors.ErrCanNotUnmarshal
	}

	return userCart, err
}

func (rr *RedisRepository) AddCart(userId uint64, userCart *models.Cart) error {
	key := rr.GetNewKey(userId)

	data, err := json.Marshal(userCart)
	if err != nil {
		return errors.ErrCanNotMarshal
	}

	_, err = redis.String(rr.conn.Do("SET", key, data))
	if err != nil {
		return errors.ErrDBInternalError
	}
	return nil
}
