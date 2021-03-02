package errors

import (
	"fmt"
	"net/http"
)

type Error struct {
	HttpCode    int    `json:"-"`
	Message     string `json:"message"`
	UserMessage string `json:"user_message"`
}

func (err Error) Error() string {
	return fmt.Sprintf("error with code %d happened %s", err.HttpCode, err.Message)
}

var (
	ErrUserUnauthorized error = Error{
		HttpCode: http.StatusUnauthorized,
		Message:  "user is unauthorized",
	}
	ErrInternalError error = Error{
		HttpCode: http.StatusInternalServerError,
		Message:  "something went wrong",
	}
)
