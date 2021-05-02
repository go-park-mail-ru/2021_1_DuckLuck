package hasher

import (
	"bytes"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

type Settings struct {
	Times   uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen int
}

func GenerateHashFromSecret(secret string, settings *Settings) ([]byte, error) {
	salt := make([]byte, settings.SaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	hashedPassword := argon2.IDKey([]byte(secret), salt, settings.Times, settings.Memory,
		settings.Threads, settings.KeyLen)
	return append(salt, hashedPassword...), nil
}

func CompareHashAndSecret(hash []byte, secret string, settings *Settings) bool {
	salt := hash[0:settings.SaltLen]
	hashedPassword := argon2.IDKey([]byte(secret), salt, settings.Times, settings.Memory,
		settings.Threads, settings.KeyLen)
	return bytes.Equal(hashedPassword, hash[settings.SaltLen:])
}
