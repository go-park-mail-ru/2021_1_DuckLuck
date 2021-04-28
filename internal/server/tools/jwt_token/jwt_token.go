package jwt_token

import (
	"encoding/base64"
	"os"
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/dgrijalva/jwt-go"
)

var (
	secretKey = os.Getenv("JWT_TOKEN")
)

type JwtToken struct {
	Value   []byte
	Expires time.Time
	jwt.StandardClaims
}

func CreateJwtToken(value []byte, expire time.Time) (string, error) {
	claims := JwtToken{
		Value:   value,
		Expires: expire,
		StandardClaims: jwt.StandardClaims{
			Issuer: "server_api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	resultToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.ErrServerSystem
	}

	return base64.StdEncoding.EncodeToString([]byte(resultToken)), nil
}

func ParseJwtToken(value string, claims *JwtToken) (*jwt.Token, error) {
	decodingValue, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return nil, errors.ErrServerSystem
	}

	jwtToken, err := jwt.ParseWithClaims(string(decodingValue), claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.ErrIncorrectJwtToken
			}
			return []byte(secretKey), nil
		})

	if err != nil {
		return nil, errors.ErrIncorrectJwtToken
	}

	return jwtToken, nil
}
