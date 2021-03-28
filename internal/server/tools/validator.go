package tools

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

func ValidateStruct(data interface{}) error {
	if _, err := govalidator.ValidateStruct(data); err != nil {
		return errors.New("incorrect field: " + err.Error())
	}

	return nil
}
