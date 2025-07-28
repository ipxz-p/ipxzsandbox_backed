package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("password", PasswordValidator)
}

var fieldNames = map[string]string{
	"ID":       "ID",
	"Name":     "Name",
	"Email":    "Email",
	"Password": "Password",
}

var validationMessages = map[string]string{
	"required": "%s is required",
	"min":      "%s must be at least %s characters long",
	"max":      "%s must not exceed %s characters",
	"email":    "%s must be a valid email address",
	"password": "%s must contain at least one uppercase letter, one lowercase letter, one digit, and one special character",
	"numeric":  "%s must be a number",
	"alpha":    "%s must contain only alphabetic characters",
	"alphanum": "%s must contain only alphanumeric characters",
	"url":      "%s must be a valid URL",
	"uuid":     "%s must be a valid UUID",
	"datetime": "%s must be a valid datetime",
	"gte":      "%s must be greater than or equal to %s",
	"lte":      "%s must be less than or equal to %s",
	"gt":       "%s must be greater than %s",
	"lt":       "%s must be less than %s",
	"eq":       "%s must be equal to %s",
	"ne":       "%s must not be equal to %s",
	"oneof":    "%s must be one of the allowed values",
	"unique":   "%s must be unique",
	"len":      "%s must be exactly %s characters long",
	"boolean":  "%s must be a boolean value",
	"json":     "%s must be valid JSON",
	"base64":   "%s must be valid base64 encoding",
	"hexcolor": "%s must be a valid hex color",
	"rgb":      "%s must be a valid RGB color",
	"rgba":     "%s must be a valid RGBA color",
	"hsl":      "%s must be a valid HSL color",
	"hsla":     "%s must be a valid HSLA color",
	"e164":     "%s must be a valid E.164 phone number format",
	"isbn":     "%s must be a valid ISBN",
	"isbn10":   "%s must be a valid ISBN-10",
	"isbn13":   "%s must be a valid ISBN-13",
	"credit_card": "%s must be a valid credit card number",
}

func getFieldName(fieldName string) string {
	if displayName, exists := fieldNames[fieldName]; exists {
		return displayName
	}
	return fieldName
}

func getValidationMessage(fieldError validator.FieldError) string {
	fieldName := getFieldName(fieldError.Field())
	tag := fieldError.Tag()
	param := fieldError.Param()

	if template, exists := validationMessages[tag]; exists {
		if param != "" {
			return fmt.Sprintf(template, fieldName, param)
		}
		return fmt.Sprintf(template, fieldName)
	}

	return fmt.Sprintf("%s failed validation on '%s' rule", fieldName, tag)
}

func TranslateValidationError(err error) map[string]string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		errsMap := make(map[string]string)
		for _, e := range errs {
			errsMap[e.Field()] = getValidationMessage(e)
		}
		return errsMap
	}
	return map[string]string{"error": err.Error()}
}