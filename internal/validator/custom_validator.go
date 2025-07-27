package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("password", PasswordValidator)
}

func TranslateValidationError(err error) map[string]string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		errsMap := make(map[string]string)
		for _, e := range errs {
			errsMap[e.Field()] = fmt.Sprintf("Field '%s' failed on '%s'", e.Field(), e.Tag())
		}
		return errsMap
	}
	return map[string]string{"error": err.Error()}
}