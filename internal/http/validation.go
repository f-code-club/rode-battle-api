package http

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func Init() (*validator.Validate, error) {
	validate := validator.New()
	if err := validate.RegisterValidation("strong_password", ValidateStrongPassword); err != nil {
		return nil, fmt.Errorf("register strong_password: %w", err)
	}
	if err := validate.RegisterValidation("name", ValidateName); err != nil {
		return nil, fmt.Errorf("register name: %w", err)
	}
	return validate, nil
}

func ValidateStrongPassword(fl validator.FieldLevel) bool {
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
