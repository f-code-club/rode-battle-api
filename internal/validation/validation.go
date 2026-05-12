package validation

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidPassword = errors.New(
		"password must contain uppercase, lowercase, number, and special character",
	)

	ErrInvalidName = errors.New(
		"name must contain only letters and spaces",
	)
)

var Validate *validator.Validate

func ValidateStrongPassword(fl validator.FieldLevel) bool {
	fmt.Println("CUSTOM PASSWORD VALIDATOR CALLED")

	password := fl.Field().String()

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true

		case unicode.IsLower(ch):
			hasLower = true

		case unicode.IsDigit(ch):
			hasNumber = true

		case strings.ContainsRune(
			`!@#$%^&*()_+-=[]{};':"\|,.<>/?`,
			ch,
		):
			hasSpecial = true
		}
	}

	return hasUpper &&
		hasLower &&
		hasNumber &&
		hasSpecial
}

func ValidateName(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	for _, c := range s {
		if !unicode.IsLetter(c) && c != ' ' {
			return false
		}
	}

	return true
}

func MapValidationError(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	for _, e := range validationErrors {

		switch e.Tag() {

		case "strong_password":
			return ErrInvalidPassword

		case "human_name":
			return ErrInvalidName
		}
	}

	return err
}
