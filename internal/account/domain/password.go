package domain

import (
	"errors"
	"unicode"
)

type Password string

func NewPassword(password string) (Password, error) {
	if len(password) < 8 {
		return "", errors.New("invalid password: must be at least 8 characters long")
	}
	var (
		hasUpper bool
		hasLower bool
		hasDigit bool
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}
	if !hasUpper {
		return "", errors.New("invalid password: must contain at least one uppercase letter")
	}
	if !hasLower {
		return "", errors.New("invalid password: must contain at least one lowercase letter")
	}
	if !hasDigit {
		return "", errors.New("invalid password: must contain at least one digit")
	}
	return Password(password), nil
}
