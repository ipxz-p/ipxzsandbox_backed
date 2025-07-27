package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		upper = regexp.MustCompile(`[A-Z]`)
		lower = regexp.MustCompile(`[a-z]`)
		special = regexp.MustCompile(`[\W_]`)
	)

	return upper.MatchString(password) &&
		lower.MatchString(password) &&
		special.MatchString(password)
}
