package errors

import (
	"fmt"
)

type Error struct {
	Message string `json:"error"`
}

func (err Error) Error() string {
	return fmt.Sprintf("error: happened %s", err.Message)
}

func CreateError(err error) error {
	if _, ok := err.(Error); ok {
		return err
	}

	return Error{Message: err.Error()}
}

var (
	ErrInternalError error = Error{
		Message: "something went wrong",
	}
	ErrDBInternalError error = Error{
		Message: "internal db error",
	}
)
