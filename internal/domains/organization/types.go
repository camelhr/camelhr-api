package organization

import (
	"errors"
	"time"

	"github.com/camelhr/camelhr-api/internal/base"
)

var ErrOrganizationDisabled = errors.New("organization is disabled")

// Organization represents an organization.
type Organization struct {
	// ID is the unique identifier of the organization.
	ID int64 `db:"organization_id"`

	// Subdomain is the subdomain of the organization.
	Subdomain string `db:"subdomain"`

	// Name is the name of the organization.
	Name string `db:"name"`

	// SuspendedAt is the timestamp when the organization was suspended.
	SuspendedAt *time.Time `db:"suspended_at"`

	// DisabledAt is the timestamp when the organization was disabled.
	DisabledAt *time.Time `db:"disabled_at"`

	// Comment represents any additional information about the organization's current state.
	Comment *string `db:"comment"`

	base.Timestamps
}

// Request represents a http request to create or update an organization.
type Request struct {
	Name string `json:"name" validate:"required,ascii,max=60"`
}

// Response represents a response of an http response organization.
type Response struct {
	ID          int64      `json:"id"`
	Subdomain   string     `json:"subdomain"`
	Name        string     `json:"name"`
	SuspendedAt *time.Time `json:"suspended_at"`
	DisabledAt  *time.Time `json:"disabled_at"`
	base.Timestamps
}
