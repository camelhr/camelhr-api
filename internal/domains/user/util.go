package user

import (
	"errors"
	"regexp"
)

const (
	emailRegexString = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	passMinLength    = 8
	passMaxLength    = 32
)

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(emailRegexString)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}

	// validate password min length
	if len(password) < passMinLength {
		return errors.New("password must be at least 8 characters")
	}
	// validate password max length
	if len(password) > passMaxLength {
		return errors.New("password must be at most 32 characters")
	}

	// at least one uppercase letter
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	if !uppercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	// at least one lowercase letter
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	if !lowercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	// at least one number
	numberRegex := regexp.MustCompile(`[0-9]`)
	if !numberRegex.MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	// at least one special character
	specialCharacterRegex := regexp.MustCompile(`[^a-zA-Z0-9]`)
	if !specialCharacterRegex.MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	// no whitespace
	whitespaceRegex := regexp.MustCompile(`\s`)
	if whitespaceRegex.MatchString(password) {
		return errors.New("password must not contain whitespace")
	}

	return nil
}
