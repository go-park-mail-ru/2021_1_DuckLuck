package tools

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/asaskevich/govalidator"
)

func ValidateStruct(data interface{}) error {
	if _, err := govalidator.ValidateStruct(data); err != nil {
		return errors.ErrInvalidData
	}

	return nil
}
