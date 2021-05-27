package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

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
	subscribes := models.Subscribes{
		Credentials: nil,
	}
	userId := uint64(2)

	t.Run("AddSubscribeUser_success", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		err := repo.AddSubscribeUser(userId, &subscribes)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("AddSubscribeUser_internal_db_error", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)
		client.Close()

		err := repo.AddSubscribeUser(userId, &subscribes)
		assert.Equal(t, err, errors.ErrDBInternalError, "unexpected error")
	})
}

func TestRedisRepository_SelectCredentialsByUserId(t *testing.T) {
	subscribes := models.Subscribes{
		Credentials: nil,
	}
	userId := uint64(2)
	incorrectUserId := uint64(4)
	data, _ := json.Marshal(subscribes)

	t.Run("SelectCredentialsByUserId_success", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		client.Set(context.Background(), fmt.Sprintf("notification:%d", userId), data, 100000)
		savedCart, err := repo.SelectCredentialsByUserId(userId)
		assert.Nil(t, err, "unexpected error")
		assert.Equal(t, *savedCart, subscribes, "incorrect id")
	})

	t.Run("SelectCredentialsByUserId_can't_parse_data", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		client.Set(context.Background(), fmt.Sprintf("notification:%d", incorrectUserId),
			"failureValue", 100000)
		_, err := repo.SelectCredentialsByUserId(incorrectUserId)
		assert.Equal(t, err, errors.ErrCanNotUnmarshal, "expected error")
	})

	t.Run("SelectCredentialsByUserId_session_not_found", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		_, err := repo.SelectCredentialsByUserId(323)
		assert.Equal(t, err, errors.ErrSessionNotFound, "expected error")
	})
}

func TestRedisRepository_DeleteSubscribeUser(t *testing.T) {
	subscribes := models.Subscribes{
		Credentials: nil,
	}
	userId := uint64(2)
	data, _ := json.Marshal(subscribes)

	t.Run("DeleteSubscribeUser_success", func(t *testing.T) {
		client := newMockRedis()
		repo := NewSessionRedisRepository(client)

		client.Set(context.Background(), fmt.Sprintf("notification:%d", userId), data, 100000)
		err := repo.DeleteSubscribeUser(userId)
		assert.Nil(t, err, "unexpected error")
	})

	t.Run("DeleteSubscribeUser_internal_db_error", func(t *testing.T) {
		client := newMockRedis()
		repo := NewSessionRedisRepository(client)

		client.Close()
		err := repo.DeleteSubscribeUser(userId)
		assert.Equal(t, err, errors.ErrDBInternalError, "unexpected error")
	})
}
