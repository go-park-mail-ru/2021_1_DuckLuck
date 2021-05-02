package errors

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/errors"
)

func CreateError(err error) error {
	return errors.CreateError(err)
}

var (
	ErrInternalError error = errors.Error{
		Message: "something went wrong",
	}
	ErrSessionNotFound error = errors.Error{
		Message: "session not found",
	}
	ErrDBInternalError error = errors.Error{
		Message: "internal db error",
	}
)
