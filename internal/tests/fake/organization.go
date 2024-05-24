package fake

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
)

// FakeOrganization is a fake organization for testing.
// It embeds the organization.Organization struct to inherit its fields.
// This is useful when you want to add custom fields or methods to the fake organization.
type FakeOrganization struct {
	organization.Organization
}

// OrganizationOption is a function that modifies an organization's default values.
type OrganizationOption func(*FakeOrganization) (*FakeOrganization, error)

// OrganizationSubdomain sets/overrides the default subdomain of an organization.
func OrganizationSubdomain(subdomain string) OrganizationOption {
	return func(o *FakeOrganization) (*FakeOrganization, error) {
		o.Subdomain = subdomain
		return o, nil
	}
}

// OrganizationName sets/overrides the default name of an organization.
func OrganizationName(name string) OrganizationOption {
	return func(o *FakeOrganization) (*FakeOrganization, error) {
		o.Name = name
		return o, nil
	}
}

// OrganizationDeleted sets deleted_at to current timestamp.
func OrganizationDeleted() OrganizationOption {
	return func(o *FakeOrganization) (*FakeOrganization, error) {
		now := time.Now().UTC()
		o.DeletedAt = &now

		return o, nil
	}
}

// OrganizationSuspended sets suspended_at to current timestamp.
func OrganizationSuspended() OrganizationOption {
	return func(o *FakeOrganization) (*FakeOrganization, error) {
		now := time.Now().UTC()
		o.SuspendedAt = &now

		return o, nil
	}
}

// OrganizationBlacklisted sets blacklisted_at to current timestamp.
func OrganizationBlacklisted() OrganizationOption {
	return func(o *FakeOrganization) (*FakeOrganization, error) {
		now := time.Now().UTC()
		o.BlacklistedAt = &now

		return o, nil
	}
}

// NewOrganization creates a fake organization for testing.
func NewOrganization(db database.Database, options ...OrganizationOption) *FakeOrganization {
	org := &FakeOrganization{}
	org.setDefaults()

	var err error
	for _, fn := range options {
		org, err = fn(org)
		if err != nil {
			panic(err)
		}
	}

	if err := org.persist(db); err != nil {
		panic(err)
	}

	return org
}

// setDefaults sets the default values of a fake organization.
//
//nolint:mnd // generates random values
func (o *FakeOrganization) setDefaults() {
	o.Subdomain = gofakeit.LetterN(uint(gofakeit.Number(1, 30)))
	o.Name = fmt.Sprint(gofakeit.LetterN(8), " ", gofakeit.Company())
	o.CreatedAt = time.Now().UTC()
	o.UpdatedAt = o.CreatedAt
}

func (o *FakeOrganization) persist(db database.Database) error {
	insertQuery := `INSERT INTO organizations
			(subdomain, name, suspended_at, blacklisted_at, comment, created_at, updated_at, deleted_at) VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING *`

	return db.Exec(context.Background(), o, insertQuery, o.Subdomain, o.Name, o.SuspendedAt, o.BlacklistedAt,
		o.Comment, o.CreatedAt, o.UpdatedAt, o.DeletedAt)
}

// IsDeleted returns deleted status of the organization by querying the database.
func (o *FakeOrganization) IsDeleted(db database.Database) bool {
	var isDeleted bool

	query := "SELECT deleted_at IS NOT NULL FROM organizations WHERE organization_id = $1"
	err := db.Get(context.Background(), &isDeleted, query, o.ID)

	if database.SuppressNoRowsError(err) != nil {
		panic(err)
	}

	return isDeleted
}

// IsSuspended returns suspended status of the organization by querying the database.
func (o *FakeOrganization) IsSuspended(db database.Database) bool {
	var isSuspended bool

	query := "SELECT suspended_at IS NOT NULL FROM organizations WHERE organization_id = $1"
	err := db.Get(context.Background(), &isSuspended, query, o.ID)

	if database.SuppressNoRowsError(err) != nil {
		panic(err)
	}

	return isSuspended
}

// IsBlacklisted returns blacklisted status of the organization by querying the database.
func (o *FakeOrganization) IsBlacklisted(db database.Database) bool {
	var isBlacklisted bool

	query := "SELECT blacklisted_at IS NOT NULL FROM organizations WHERE organization_id = $1"
	err := db.Get(context.Background(), &isBlacklisted, query, o.ID)

	if database.SuppressNoRowsError(err) != nil {
		panic(err)
	}

	return isBlacklisted
}

// AddUser creates a fake user under the organization.
func (o *FakeOrganization) AddUser(db database.Database, userOptions ...UserOption) *FakeUser {
	user := NewUser(db, o.ID, userOptions...)
	return user
}

// FetchLatest fetches and returns the latest version of organization by querying the database.
func (o *FakeOrganization) FetchLatest(db database.Database) *FakeOrganization {
	fakeOrg := &FakeOrganization{}

	query := `
			SELECT
				organization_id,
				subdomain,
				name,
				suspended_at,
				blacklisted_at,
				comment,
				created_at,
				updated_at,
				deleted_at
			FROM organizations 
			WHERE organization_id = $1
			`

	if err := db.Get(context.Background(), fakeOrg, query, o.ID); err != nil {
		panic(err)
	}

	return fakeOrg
}
