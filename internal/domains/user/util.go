package user

import (
	"regexp"

	"github.com/camelhr/camelhr-api/internal/base"
)

const (
	emailRegexString = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	passMinLength    = 8
	passMaxLength    = 32
)

// ValidateEmail validates the email string.
func ValidateEmail(email string) error {
	if email == "" {
		return base.NewInputValidationError("email is required")
	}

	emailRegex := regexp.MustCompile(emailRegexString)
	if !emailRegex.MatchString(email) {
		return base.NewInputValidationError("email must be a valid email address")
	}

	return nil
}

// ValidatePassword validates the password string.
func ValidatePassword(password string) error {
	if password == "" {
		return base.NewInputValidationError("password is required")
	}

	// validate password min length
	if len(password) < passMinLength {
		return base.NewInputValidationError("password must be at least 8 characters in length")
	}
	// validate password max length
	if len(password) > passMaxLength {
		return base.NewInputValidationError("password must be at most 32 characters")
	}

	// at least one uppercase letter
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	if !uppercaseRegex.MatchString(password) {
		return base.NewInputValidationError("password must contain at least one uppercase letter")
	}

	// at least one lowercase letter
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	if !lowercaseRegex.MatchString(password) {
		return base.NewInputValidationError("password must contain at least one lowercase letter")
	}

	// at least one number
	numberRegex := regexp.MustCompile(`[0-9]`)
	if !numberRegex.MatchString(password) {
		return base.NewInputValidationError("password must contain at least one number")
	}

	// at least one special character
	specialCharacterRegex := regexp.MustCompile(`[^a-zA-Z0-9]`)
	if !specialCharacterRegex.MatchString(password) {
		return base.NewInputValidationError("password must contain at least one special character")
	}

	// no whitespace
	whitespaceRegex := regexp.MustCompile(`\s`)
	if whitespaceRegex.MatchString(password) {
		return base.NewInputValidationError("password must not contain whitespace")
	}

	return nil
}

// ValidateComment validates the comment string.
func ValidateComment(comment string) error {
	const allowedMaxLength = 255

	// validate that comment is not empty
	if comment == "" {
		return base.NewInputValidationError("comment is required")
	}

	// validate that comment length does not exceed allowedMaxLength
	if len(comment) > allowedMaxLength {
		return base.NewInputValidationError("comment must be a maximum of 255 characters in length")
	}

	return nil
}
