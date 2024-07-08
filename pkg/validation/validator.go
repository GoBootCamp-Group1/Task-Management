package validation

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

type Validator struct {
	*validator.Validate
}

var validate *validator.Validate

func NewValidator() *Validator {
	if validate == nil {
		validate = validator.New(validator.WithRequiredStructEnabled())
		validate.RegisterValidation("password", validatePassword)
	}
	return &Validator{Validate: validate}
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	minLen := 8

	if len(password) < minLen {
		return false
	}

	upperCaseRegex := regexp.MustCompile(`[A-Z]`)
	specialCharRegex := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)

	if !upperCaseRegex.MatchString(password) {
		return false
	}
	if !specialCharRegex.MatchString(password) {
		return false
	}

	return true
}
