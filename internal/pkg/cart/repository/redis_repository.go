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

func (r *RedisRepository) getNewKey(value uint64) string {
	return fmt.Sprintf("cart:%d", value)
}

// Delete user cart
func (r *RedisRepository) DeleteCart(userId uint64) error {
	key := r.getNewKey(userId)

	result, err := redis.String(r.conn.Do("DEL", key))
	if err != nil || result != "OK" {
		return errors.ErrDBInternalError
	}
	return nil
}

// Select user cart by id
func (r *RedisRepository) SelectCartById(userId uint64) (*models.Cart, error) {
	userCart := &models.Cart{}
	key := r.getNewKey(userId)

	data, err := redis.Bytes(r.conn.Do("GET", key))
	if err != nil {
		return nil, errors.ErrCartNotFound
	}

	if err = json.Unmarshal(data, &userCart); err != nil {
		return nil, errors.ErrCanNotUnmarshal
	}

	return userCart, err
}

// Add new user cart
func (r *RedisRepository) AddCart(userId uint64, userCart *models.Cart) error {
	key := r.getNewKey(userId)

	data, err := json.Marshal(userCart)
	if err != nil {
		return errors.ErrCanNotMarshal
	}

	result, err := redis.String(r.conn.Do("SET", key, data))
	if err != nil || result != "OK" {
		return errors.ErrDBInternalError
	}
	return nil
}
