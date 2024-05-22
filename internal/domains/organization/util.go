package organization

import (
	"errors"
	"regexp"
)

// ValidateSubdomain validates the subdomain string.
func ValidateSubdomain(subdomain string) error {
	const allowedMaxLength = 30

	// validate that subdomain is not empty
	if subdomain == "" {
		return errors.New("subdomain is a required")
	}

	// validate that subdomain length does not exceed allowedMaxLength
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

// ValidateOrgName validates the organization name string.
func ValidateOrgName(orgName string) error {
	const allowedMaxLength = 60

	// validate that orgName is not empty
	if orgName == "" {
		return errors.New("organization name is required")
	}

	// validate that orgName length does not exceed allowedMaxLength
	if len(orgName) > allowedMaxLength {
		return errors.New("organization name must be a maximum of 60 characters in length")
	}

	// validate that orgName is ascii
	match, err := regexp.MatchString(`^[\x00-\x7F]+$`, orgName)
	if err != nil || !match {
		return errors.New("organization name can only contain ascii characters")
	}

	return nil
}
