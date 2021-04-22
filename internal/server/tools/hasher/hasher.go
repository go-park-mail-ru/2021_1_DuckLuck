package hasher

import (
	"bytes"
	"crypto/rand"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"golang.org/x/crypto/argon2"
)

type Settings struct {
	times   uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen int
}

var (
	passwordSettings = &Settings{
		times:   1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  64,
		saltLen: 8,
	}
)

func GenerateHashFromPassword(password string) ([]byte, error) {
	return generateHashFromSecret(password, passwordSettings)
}

func CompareHashAndPassword(hash []byte, password string) bool {
	return compareHashAndSecret(hash, password, passwordSettings)
}

func generateHashFromSecret(secret string, settings *Settings) ([]byte, error) {
	salt := make([]byte, settings.saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, errors.ErrServerSystem
	}

	hashedPassword := argon2.IDKey([]byte(secret), salt, settings.times, settings.memory,
		settings.threads, settings.keyLen)
	return append(salt, hashedPassword...), nil
}

func compareHashAndSecret(hash []byte, secret string, settings *Settings) bool {
	salt := hash[0:settings.saltLen]
	hashedPassword := argon2.IDKey([]byte(secret), salt, settings.times, settings.memory,
		settings.threads, settings.keyLen)
	return bytes.Equal(hashedPassword, hash[settings.saltLen:])
}
