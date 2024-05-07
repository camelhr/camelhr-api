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
		now := time.Now()
		o.DeletedAt = &now

		return o, nil
	}
}

// OrganizationSuspended sets suspended_at to current timestamp.
func OrganizationSuspended() OrganizationOption {
	return func(o *FakeOrganization) (*FakeOrganization, error) {
		now := time.Now()
		o.SuspendedAt = &now

		return o, nil
	}
}

// OrganizationBlacklisted sets blacklisted_at to current timestamp.
func OrganizationBlacklisted() OrganizationOption {
	return func(o *FakeOrganization) (*FakeOrganization, error) {
		now := time.Now()
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

	if err := persist(db, org); err != nil {
		panic(err)
	}

	return org
}

//nolint:gomnd // generate random values
func (o *FakeOrganization) setDefaults() {
	o.Subdomain = gofakeit.LetterN(uint(gofakeit.Number(1, 30)))
	o.Name = fmt.Sprint(gofakeit.LetterN(8), " ", gofakeit.Company())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = time.Now().UTC()
	}

	if o.UpdatedAt.IsZero() {
		o.UpdatedAt = o.CreatedAt
	}
}

func persist(db database.Database, o *FakeOrganization) error {
	insertQuery := `INSERT INTO organizations
			(subdomain, name, suspended_at, blacklisted_at, comment, created_at, updated_at, deleted_at) VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING *`

	if err := db.Exec(context.TODO(), o, insertQuery, o.Subdomain, o.Name, o.SuspendedAt, o.BlacklistedAt,
		o.Comment, o.CreatedAt, o.UpdatedAt, o.DeletedAt); err != nil {
		return err
	}

	return nil
}

func (o *FakeOrganization) IsDeleted(db database.Database) bool {
	var isDeleted bool

	query := "SELECT deleted_at IS NOT NULL FROM organizations WHERE organization_id = $1"
	err := db.Get(context.TODO(), &isDeleted, query, o.ID)

	if database.SuppressNoRowsError(err) != nil {
		panic(err)
	}

	return isDeleted
}

// IsSuspended returns suspended status of the organization by querying the database.
func (o *FakeOrganization) IsSuspended(db database.Database) bool {
	var isSuspended bool

	query := "SELECT suspended_at IS NOT NULL FROM organizations WHERE organization_id = $1"
	err := db.Get(context.TODO(), &isSuspended, query, o.ID)

	if database.SuppressNoRowsError(err) != nil {
		panic(err)
	}

	return isSuspended
}

// IsBlacklisted returns blacklisted status of the organization by querying the database.
func (o *FakeOrganization) IsBlacklisted(db database.Database) bool {
	var isBlacklisted bool

	query := "SELECT blacklisted_at IS NOT NULL FROM organizations WHERE organization_id = $1"
	err := db.Get(context.TODO(), &isBlacklisted, query, o.ID)

	if database.SuppressNoRowsError(err) != nil {
		panic(err)
	}

	return isBlacklisted
}

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

	if err := db.Get(context.TODO(), fakeOrg, query, o.ID); err != nil {
		panic(err)
	}

	return fakeOrg
}
