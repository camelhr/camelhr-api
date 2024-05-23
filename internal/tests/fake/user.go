package fake

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"golang.org/x/crypto/bcrypt"
)

// FakeUser is a fake user for testing.
// It embeds the user.User struct to inherit its fields.
// This is useful when you want to add custom fields or methods to the fake user.
type FakeUser struct {
	user.User
	setAPIToken bool
}

// UserOption is a function that modifies a user's default values.
type UserOption func(*FakeUser) (*FakeUser, error)

// UserEmail sets/overrides the default email of a user.
func UserEmail(email string) UserOption {
	return func(u *FakeUser) (*FakeUser, error) {
		u.Email = email
		return u, nil
	}
}

// UserPassword sets/overrides the default password of a user.
func UserPassword(password string) UserOption {
	return func(u *FakeUser) (*FakeUser, error) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		u.PasswordHash = string(hashedPassword)

		return u, nil
	}
}

// UserIsOwner sets the user as the owner of the organization.
func UserIsOwner() UserOption {
	return func(u *FakeUser) (*FakeUser, error) {
		u.IsOwner = true
		return u, nil
	}
}

// UserEmailNotVerified sets is_email_verified to false.
func UserEmailNotVerified() UserOption {
	return func(u *FakeUser) (*FakeUser, error) {
		u.IsEmailVerified = false
		return u, nil
	}
}

// UserWithoutToken sets the API token to null.
func UserWithoutToken() UserOption {
	return func(u *FakeUser) (*FakeUser, error) {
		u.setAPIToken = false
		return u, nil
	}
}

// UserDisabled sets disabled_at to current timestamp.
func UserDisabled() UserOption {
	return func(u *FakeUser) (*FakeUser, error) {
		now := time.Now().UTC()
		u.DisabledAt = &now

		return u, nil
	}
}

// UserDeleted sets deleted_at to current timestamp.
func UserDeleted() UserOption {
	return func(u *FakeUser) (*FakeUser, error) {
		now := time.Now().UTC()
		u.DeletedAt = &now

		return u, nil
	}
}

// NewUser creates a fake user for testing.
func NewUser(db database.Database, orgID int64, options ...UserOption) *FakeUser {
	u := &FakeUser{}
	u.OrganizationID = orgID
	u.setDefaults()

	var err error
	for _, fn := range options {
		u, err = fn(u)
		if err != nil {
			panic(err)
		}
	}

	if err := u.persist(db); err != nil {
		panic(err)
	}

	return u
}

// IsDisabled returns disabled status of the user by querying the database.
func (u *FakeUser) IsDisabled(db database.Database) bool {
	u = u.FetchLatest(db)
	return u.DisabledAt != nil
}

// IsDeleted returns deleted status of the user by querying the database.
func (u *FakeUser) IsDeleted(db database.Database) bool {
	u = u.FetchLatest(db)
	return u.DeletedAt != nil
}

// setDefaults sets the default values of a fake user.
//
//nolint:gomnd // generate random values
func (u *FakeUser) setDefaults() {
	u.Email = gofakeit.Email()

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(gofakeit.Password(true, true, true, true, false, 12)),
		bcrypt.DefaultCost,
	)
	if err != nil {
		panic(err)
	}

	u.PasswordHash = string(hashedPassword)
	u.IsOwner = false
	u.IsEmailVerified = true
	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = u.CreatedAt
	u.setAPIToken = true
}

// persist saves the fake user to the database.
func (u *FakeUser) persist(db database.Database) error {
	insertQuery := `INSERT INTO users
	(organization_id, email, password_hash, api_token, is_owner, is_email_verified,
		disabled_at, comment, created_at, updated_at, deleted_at) VALUES
	($1, $2, $3, CASE WHEN $4 THEN random_token() ELSE null END, $5, $6, $7, $8, $9, $10, $11)
	RETURNING *`

	return db.Exec(context.Background(), u, insertQuery,
		u.OrganizationID, u.Email, u.PasswordHash, u.setAPIToken, u.IsOwner, u.IsEmailVerified,
		u.DisabledAt, u.Comment, u.CreatedAt, u.UpdatedAt, u.DeletedAt)
}

// FetchLatest fetches and returns the latest version of user by querying the database.
func (u *FakeUser) FetchLatest(db database.Database) *FakeUser {
	fakeUser := &FakeUser{}

	query := `
			SELECT
				user_id,
				organization_id,
				email,
				password_hash,
				api_token,
				is_owner,
				is_email_verified,
				disabled_at,
				comment,
				created_at,
				updated_at,
				deleted_at
			FROM users
			WHERE user_id = $1
			`

	if err := db.Get(context.Background(), fakeUser, query, u.ID); err != nil {
		panic(err)
	}

	return fakeUser
}
