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
		Message: "session not found",
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
	ErrFileNotRead error = Error{
		Message: "can't read file",
	}
	ErrIncorrectFileType error = Error{
		Message: "incorrect file type",
	}
	ErrProductNotFound error = Error{
		Message: "product not found",
	}
	ErrProductsIsEmpty error = Error{
		Message: "list of products is empty",
	}
	ErrIncorrectPaginator error = Error{
		Message: "incorrect params of pagination",
	}
	ErrBadRequest error = Error{
		Message: "incorrect request",
	}
	ErrCanNotUnmarshal error = Error{
		Message: "can't unmarshal",
	}
	ErrCanNotMarshal error = Error{
		Message: "can't marshal",
	}
	ErrDBInternalError error = Error{
		Message: "internal db error",
	}
	ErrDBFailedConnection error = Error{
		Message: "can't connect to db",
	}
	ErrHashFunctionFailed error = Error{
		Message: "can't get hash of data",
	}
	ErrCartNotFound error = Error{
		Message: "user cart not found",
	}
	ErrProductNotFoundInCart error = Error{
		Message: "product not found in cart",
	}
	ErrInvalidData error = Error{
		Message: "invalid data",
	}
)
