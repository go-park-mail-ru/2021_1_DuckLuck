package repository

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/pkg/models"

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

func TestRedisRepository_AddSession(t *testing.T) {
	sess := models.Session{
		Value: "test",
		UserId: uint64(33),
	}

	t.Run("AddSession_success", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		err := repo.AddSession(&sess)
		assert.Nil(t, err, "unexpected error")
	})

	t.Run("AddSession_internal_db_error", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		client.Close()
		err := repo.AddSession(&sess)
		assert.Equal(t, err, errors.ErrDBInternalError, "unexpected error")
	})
}

func TestRedisRepository_SelectUserIdBySession(t *testing.T) {
	sessionValue := "test"
	incorrectValue := "none"
	failureValue := "five"
	sess := models.Session{
		Value: sessionValue,
		UserId: uint64(33),
	}
	data := fmt.Sprintf("%d", sess.UserId)

	t.Run("SelectUserIdBySession_success", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		client.Set(context.Background(), "session:"+sessionValue, data, 100000)
		userId, err := repo.SelectUserIdBySession(sessionValue)
		assert.Nil(t, err, "unexpected error")
		assert.Equal(t, sess.UserId, userId, "incorrect id")
	})

	t.Run("SelectUserIdBySession_internal_db_error", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		client.Set(context.Background(), "session:"+failureValue, failureValue, 100000)
		_, err := repo.SelectUserIdBySession(failureValue)
		assert.Equal(t, err, errors.ErrInternalError, "expected error")
	})

	t.Run("SelectUserIdBySession_session_not_found", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		_, err := repo.SelectUserIdBySession(incorrectValue)
		assert.Equal(t, err, errors.ErrSessionNotFound, "expected error")
	})
}

func TestRedisRepository_DeleteSessionByValue(t *testing.T) {
	sessionValue := "test"
	sess := models.Session{
		Value: sessionValue,
		UserId: uint64(33),
	}
	data := fmt.Sprintf("%d", sess.UserId)

	t.Run("DeleteSessionByValue_success", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		client.Set(context.Background(), "session:"+sessionValue, data, 100000)
		err := repo.DeleteSessionByValue(sessionValue)
		assert.Nil(t, err, "unexpected error")
	})

	t.Run("DeleteSessionByValue_internal_db_error", func(t *testing.T) {
		client := newMockRedis()
		defer client.Close()
		repo := NewSessionRedisRepository(client)

		client.Close()
		err := repo.DeleteSessionByValue(sessionValue)
		assert.Equal(t, err, errors.ErrDBInternalError, "unexpected error")
	})
}