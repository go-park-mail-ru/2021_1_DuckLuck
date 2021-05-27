package password_hasher

import "github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/hasher"

var (
	passwordSettings = &hasher.Settings{
		Times:   1,
		Memory:  1 * 1024,
		Threads: 1,
		KeyLen:  32,
		SaltLen: 4,
	}
)

func GenerateHashFromPassword(password string) ([]byte, error) {
	return hasher.GenerateHashFromSecret(password, passwordSettings)
}

func CompareHashAndPassword(hash []byte, password string) bool {
	return hasher.CompareHashAndSecret(hash, password, passwordSettings)
}
