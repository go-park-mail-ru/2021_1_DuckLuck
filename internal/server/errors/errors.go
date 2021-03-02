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
	ErrUserNotFound error = Error{
		HttpCode: http.StatusInternalServerError,
		Message:  "user not found",
	}
	ErrSessionNotFound error = Error{
		HttpCode: http.StatusInternalServerError,
		Message:  "something went wrong",
	}
	ErrIncorrectUserEmail error = Error{
		HttpCode: http.StatusBadRequest,
		Message:  "incorrect user email",
	}
	ErrIncorrectUserPassword error = Error{
		HttpCode: http.StatusBadRequest,
		Message:  "incorrect user password",
	}
	ErrEmailAlreadyExist error = Error{
		HttpCode: http.StatusConflict,
		Message:  "user email already exist",
	}
)
