package user

import (
	"time"

	"github.com/camelhr/camelhr-api/internal/base"
)

// User represents a user.
type User struct {
	// ID is the unique identifier of the user.
	ID int64 `db:"user_id"`

	// OrganizationID is the reference to the organization the user belongs to.
	OrganizationID int64 `db:"organization_id"`

	// Email is the email address of the user.
	Email string `db:"email"`

	// PasswordHash is the hashed password of the user.
	PasswordHash string `db:"password_hash"`

	// APIToken is the token used to authenticate the user for API requests.
	APIToken *string `db:"api_token"`

	// IsOwner represents whether the user is the owner of the organization.
	IsOwner bool `db:"is_owner"`

	// DisabledAt is the timestamp when the user was disabled.
	DisabledAt *time.Time `db:"disabled_at"`

	// Comment represents any additional information about the user's current state.
	Comment *string `db:"comment"`

	base.Timestamps
}
