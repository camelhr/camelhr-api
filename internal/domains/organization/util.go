package organization

import (
	"regexp"

	"github.com/camelhr/camelhr-api/internal/base"
)

// ValidateSubdomain validates the subdomain string.
func ValidateSubdomain(subdomain string) error {
	const allowedMaxLength = 30

	// validate that subdomain is not empty
	if subdomain == "" {
		return base.NewInputValidationError("subdomain is required")
	}

	// validate that subdomain length does not exceed allowedMaxLength
	if len(subdomain) > allowedMaxLength {
		return base.NewInputValidationError("subdomain must be a maximum of 30 characters in length")
	}

	// validate that subdomain is alphanumeric
	match, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, subdomain)
	if err != nil || !match {
		return base.NewInputValidationError("subdomain can only contain alphanumeric characters")
	}

	return nil
}

// ValidateOrgName validates the organization name string.
func ValidateOrgName(orgName string) error {
	const allowedMaxLength = 60

	// validate that orgName is not empty
	if orgName == "" {
		return base.NewInputValidationError("organization name is required")
	}

	// validate that orgName length does not exceed allowedMaxLength
	if len(orgName) > allowedMaxLength {
		return base.NewInputValidationError("organization name must be a maximum of 60 characters in length")
	}

	// validate that orgName is ascii
	match, err := regexp.MatchString(`^[\x00-\x7F]+$`, orgName)
	if err != nil || !match {
		return base.NewInputValidationError("organization name can only contain ascii characters")
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
