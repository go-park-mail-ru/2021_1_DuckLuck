package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/pkg/models"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func newMockRedis() *redis.Client {
	conn, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: conn.Addr(),
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err)
	}

	return client
}

func TestNewSessionRedisRepository(t *testing.T) {
	assert.NotEmpty(t, NewSessionRedisRepository(newMockRedis()),
		"can't create new repository")
}

func TestRedisRepository_AddCart(t *testing.T) {
	cart := models.Cart{
		Products: nil,
	}
	userId := uint64(2)

	t.Run("AddCart_success", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		err := repo.AddCart(userId, &cart)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("AddCart_internal_db_error", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)
		client.Close()

		err := repo.AddCart(userId, &cart)
		assert.Equal(t, err, errors.ErrDBInternalError, "unexpected error")
	})
}

func TestRedisRepository_SelectCartById(t *testing.T) {
	cart := models.Cart{
		Products: nil,
	}
	userId := uint64(2)
	incorrectUserId := uint64(4)
	data, _ := json.Marshal(cart)

	t.Run("SelectCartById_success", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		client.Set(context.Background(), fmt.Sprintf("cart:%d", userId), data, 100000)
		savedCart, err := repo.SelectCartById(userId)
		assert.Nil(t, err, "unexpected error")
		assert.Equal(t, *savedCart, cart, "incorrect id")
	})

	t.Run("SelectCartById_can't_parse_data", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		client.Set(context.Background(), fmt.Sprintf("cart:%d", incorrectUserId),
			"failureValue", 100000)
		_, err := repo.SelectCartById(incorrectUserId)
		assert.Equal(t, err, errors.ErrCanNotUnmarshal, "expected error")
	})

	t.Run("SelectCartById_session_not_found", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		_, err := repo.SelectCartById(323)
		assert.Equal(t, err, errors.ErrCartNotFound, "expected error")
	})
}

func TestRedisRepository_DeleteCart(t *testing.T) {
	cart := models.Cart{
		Products: nil,
	}
	userId := uint64(2)
	data, _ := json.Marshal(cart)

	t.Run("DeleteCart_success", func(t *testing.T) {
		client := newMockRedis()
		repo := NewSessionRedisRepository(client)

		client.Set(context.Background(), fmt.Sprintf("cart:%d", userId), data, 100000)
		err := repo.DeleteCart(userId)
		assert.Nil(t, err, "unexpected error")
	})

	t.Run("DeleteCart_internal_db_error", func(t *testing.T) {
		client := newMockRedis()
		repo := NewSessionRedisRepository(client)

		client.Close()
		err := repo.DeleteCart(userId)
		assert.Equal(t, err, errors.ErrDBInternalError, "unexpected error")
	})
}