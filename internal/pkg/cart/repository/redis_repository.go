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

func NewSessionPostgresqlRepository(conn redis.Conn) cart.Repository {
	return &RedisRepository{
		conn: conn,
	}
}

func (rr *RedisRepository) GetNewKey(value uint64) string {
	return fmt.Sprintf("cart:%d", value)
}

func (rr *RedisRepository) AddProductPosition(userId uint64, position *models.ProductPosition) error {
	userCart := &models.Cart{}
	key := rr.GetNewKey(userId)

	data, err := redis.Bytes(rr.conn.Do("GET", key))
	if err == nil {
		if err = json.Unmarshal(data, &userCart); err != nil {
			return errors.ErrCanNotUnmarshal
		}
	}

	// If product position already exist then increment counter
	if _, ok := userCart.Products[position.ProductId]; ok {
		userCart.Products[position.ProductId].Count++
	} else {
		userCart.Products[position.ProductId] = position
	}

	return rr.AddProductsInCart(userId, userCart)
}

func (rr *RedisRepository) DeleteProductPosition(userId uint64, productId uint64) error {
	userCart, err := rr.GetProductsFromCart(userId)
	if err != nil {
		return err
	}

	// Delete cart if empty
	key := rr.GetNewKey(userId)
	if len(userCart.Products) == 1 {
		_, err = redis.String(rr.conn.Do("DEL", key))
		if err != nil {
			return errors.ErrDBInternalError
		}
		return nil
	}

	// Delete only product position
	delete(userCart.Products, productId)
	return rr.AddProductsInCart(userId, userCart)
}

func (rr *RedisRepository) UpdateProductPosition(userId uint64, position *models.ProductPosition) error {
	userCart, err := rr.GetProductsFromCart(userId)
	if err != nil {
		return err
	}

	if _, ok := userCart.Products[position.ProductId]; !ok {
		return errors.ErrProductNotFoundInCart
	}
	userCart.Products[position.ProductId] = position

	return rr.AddProductsInCart(userId, userCart)
}

func (rr *RedisRepository) GetProductsFromCart(userId uint64) (*models.Cart, error) {
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

func (rr *RedisRepository) AddProductsInCart(userId uint64, userCart *models.Cart) error {
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
