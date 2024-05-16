package organization

import (
	"errors"
	"regexp"
)

func ValidateSubdomain(subdomain string) error {
	const allowedMaxLength = 30

	// validate that subdomain is not empty
	if subdomain == "" {
		return errors.New("subdomain is a required")
	}

	// validate that subdomain length is less than or equal to 30
	if len(subdomain) > allowedMaxLength {
		return errors.New("subdomain must be a maximum of 30 characters in length")
	}

	// validate that subdomain is alphanumeric
	match, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, subdomain)
	if err != nil || !match {
		return errors.New("subdomain can only contain alphanumeric characters")
	}

	return nil
}
