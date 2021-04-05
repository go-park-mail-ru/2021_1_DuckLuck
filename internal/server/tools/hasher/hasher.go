package password_hasher

import (
	"bytes"
	"crypto/rand"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"golang.org/x/crypto/argon2"
)

const (
	times   = 1
	memory  = 64 * 1024
	threads = 4
	keyLen  = 32
	saltLen = 8
)

func GenerateHashFromPassword(password string) ([]byte, error) {
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, errors.ErrServerSystem
	}

	hashedPassword := argon2.IDKey([]byte(password), salt, times, memory, threads, keyLen)
	return append(salt, hashedPassword...), nil
}

func CompareHashAndPassword(hash []byte, password string) bool {
	salt := hash[0:saltLen]
	hashedPassword := argon2.IDKey([]byte(password), salt, times, memory, threads, keyLen)
	return bytes.Equal(hashedPassword, hash[saltLen:])
}
