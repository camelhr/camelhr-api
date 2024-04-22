package organization

import (
	"time"

	"github.com/camelhr/camelhr-api/internal/base"
)

// Organization represents an organization.
type Organization struct {
	// ID is the unique identifier of the organization.
	ID int64 `db:"organization_id"`

	// Name is the name of the organization.
	Name string `db:"name"`

	// SuspendedAt is the timestamp when the organization was suspended.
	SuspendedAt *time.Time `db:"suspended_at"`

	// BlacklistedAt is the timestamp when the organization was blacklisted.
	BlacklistedAt *time.Time `db:"blacklisted_at"`

	// Comment represents any additional information about the organization's current state.
	Comment string `db:"comment"`

	base.Timestamps
}

// Request represents a http request to create or update an organization.
type Request struct {
	Name string `json:"name" validate:"required"`
}

// Response represents a response of an http response organization.
type Response struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	SuspendedAt   *time.Time `json:"suspended_at"`
	BlacklistedAt *time.Time `json:"blacklisted_at"`
	base.Timestamps
}
