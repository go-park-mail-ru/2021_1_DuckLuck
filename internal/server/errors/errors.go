package errors

import (
	"fmt"
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
		Message: "user is unauthorized",
	}
	ErrInternalError error = Error{
		Message: "something went wrong",
	}
	ErrUserNotFound error = Error{
		Message: "user not found",
	}
	ErrSessionNotFound error = Error{
		Message: "something went wrong",
	}
	ErrIncorrectUserEmail error = Error{
		Message: "incorrect user email",
	}
	ErrIncorrectUserPassword error = Error{
		Message: "incorrect user password",
	}
	ErrEmailAlreadyExist error = Error{
		Message: "user email already exist",
	}
	ErrServerSystem error = Error{
		Message: "system error",
	}
)
