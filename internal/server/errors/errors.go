package errors

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/errors"
)

func CreateError(err error) error {
	return errors.CreateError(err)
}

var (
	ErrUserUnauthorized error = errors.Error{
		Message: "user is unauthorized",
	}
	ErrInternalError error = errors.Error{
		Message: "something went wrong",
	}
	ErrUserNotFound error = errors.Error{
		Message: "user not found",
	}
	ErrSessionNotFound error = errors.Error{
		Message: "session not found",
	}
	ErrEmailAlreadyExist error = errors.Error{
		Message: "user email already exist",
	}
	ErrServerSystem error = errors.Error{
		Message: "system error",
	}
	ErrFileNotRead error = errors.Error{
		Message: "can't read file",
	}
	ErrIncorrectFileType error = errors.Error{
		Message: "incorrect file type",
	}
	ErrProductNotFound error = errors.Error{
		Message: "product not found",
	}
	ErrCategoryNotFound error = errors.Error{
		Message: "category not found",
	}
	ErrIncorrectPaginator error = errors.Error{
		Message: "incorrect params of pagination",
	}
	ErrBadRequest error = errors.Error{
		Message: "incorrect request",
	}
	ErrCanNotUnmarshal error = errors.Error{
		Message: "can't unmarshal",
	}
	ErrCanNotMarshal error = errors.Error{
		Message: "can't marshal",
	}
	ErrDBInternalError error = errors.Error{
		Message: "internal db error",
	}
	ErrDBFailedConnection error = errors.Error{
		Message: "can't connect to db",
	}
	ErrCartNotFound error = errors.Error{
		Message: "user cart not found",
	}
	ErrProductNotFoundInCart error = errors.Error{
		Message: "product not found in cart",
	}
	ErrInvalidData error = errors.Error{
		Message: "invalid data",
	}
	ErrRequireIdNotFound error = errors.Error{
		Message: "require id not found",
	}
	ErrOpenFile error = errors.Error{
		Message: "can't open file",
	}
	ErrNotFoundCsrfToken error = errors.Error{
		Message: "csrf token not found",
	}
	ErrIncorrectJwtToken error = errors.Error{
		Message: "incorrect jwt token",
	}
	ErrS3InternalError error = errors.Error{
		Message: "can't upload file to S3",
	}
	ErrIncorrectAuthData error = errors.Error{
		Message: "incorrect auth user data",
	}
	ErrNoWriteRights error = errors.Error{
		Message: "no write rights",
	}
	ErrCanNotAddReview error = errors.Error{
		Message: "user can not add review",
	}
	ErrIncorrectSearchQuery error = errors.Error{
		Message: "incorrect search query",
	}
	ErrPromoCodeNotFound error = errors.Error{
		Message: "promo code not found",
	}
	ErrProductNotInPromo error = errors.Error{
		Message: "product does not participate in the promotion",
	}
)
