package fake_test

import (
	"testing"
	"time"

	"github.com/camelhr/camelhr-api/internal/tests/fake"
	"github.com/stretchr/testify/assert"
)

func (s *FakeTestSuite) TestFakeOrganization() {
	s.Run("should create an organization with default values", func() {
		s.T().Parallel()

		// create a fake organization with default values
		o := fake.NewOrganization(s.DB)

		// assert that the organization is created with the default values
		s.Require().NotNil(o)
		s.Empty(o.Comment)
		s.NotEmpty(o.ID)
		s.NotEmpty(o.Name)
		s.NotNil(o.CreatedAt)
		s.NotNil(o.UpdatedAt)
		s.Equal(time.UTC, o.CreatedAt.Location())
		s.Equal(time.UTC, o.UpdatedAt.Location())
		s.WithinDuration(time.Now().UTC(), o.CreatedAt, 1*time.Minute)
		s.WithinDuration(time.Now().UTC(), o.UpdatedAt, 1*time.Minute)
		s.Nil(o.DeletedAt)
		s.Nil(o.SuspendedAt)
		s.Nil(o.BlacklistedAt)
	})

	s.Run("should create an organization with custom subdomain", func() {
		s.T().Parallel()

		// create a fake organization with custom subdomain
		subdomain := "test"
		o := fake.NewOrganization(s.DB, fake.OrganizationSubdomain(subdomain))

		// assert that the organization is created with the specified subdomain
		s.Require().NotNil(o)
		s.Equal(subdomain, o.Subdomain)
	})

	s.Run("should create an organization with custom name", func() {
		s.T().Parallel()

		// create a fake organization with custom name
		name := "test organization"
		o := fake.NewOrganization(s.DB, fake.OrganizationName(name))

		// assert that the organization is created the specified name
		s.Require().NotNil(o)
		s.Equal(name, o.Name)
	})

	s.Run("should create a deleted organization", func() {
		s.T().Parallel()

		// create a deleted organization
		o := fake.NewOrganization(s.DB, fake.OrganizationDeleted())

		// assert that the organization is set deleted
		s.Require().NotNil(o)
		s.Require().NotNil(o.DeletedAt)
		s.Equal(time.UTC, o.DeletedAt.Location())
		s.WithinDuration(time.Now().UTC(), *o.DeletedAt, 1*time.Minute)
	})

	s.Run("should create a suspended organization", func() {
		s.T().Parallel()

		// create a suspended organization
		o := fake.NewOrganization(s.DB, fake.OrganizationSuspended())

		// assert that the organization is suspended
		s.Require().NotNil(o)
		s.Require().NotNil(o.SuspendedAt)
		s.Equal(time.UTC, o.SuspendedAt.Location())
		s.WithinDuration(time.Now().UTC(), *o.SuspendedAt, 1*time.Minute)
		s.Nil(o.DeletedAt)
		s.Nil(o.BlacklistedAt)
	})

	s.Run("should return true if organization is suspended", func() {
		s.T().Parallel()

		// create a suspended organization
		o := fake.NewOrganization(s.DB, fake.OrganizationSuspended())
		s.Require().NotNil(o)
		isSuspended := o.IsSuspended(s.DB)

		// assert that the organization is suspended
		s.True(isSuspended)
	})

	s.Run("should created a blacklisted organization", func() {
		s.T().Parallel()

		// create a blacklisted organization
		o := fake.NewOrganization(s.DB, fake.OrganizationBlacklisted())

		// assert that the organization is blacklisted
		s.Require().NotNil(o)
		s.Require().NotNil(o.BlacklistedAt)
		s.Equal(time.UTC, o.BlacklistedAt.Location())
		s.WithinDuration(time.Now().UTC(), *o.BlacklistedAt, 1*time.Minute)
		s.Nil(o.DeletedAt)
		s.Nil(o.SuspendedAt)
	})

	s.Run("should return true if organization is blacklisted", func() {
		s.T().Parallel()

		// create a blacklisted organization
		o := fake.NewOrganization(s.DB, fake.OrganizationBlacklisted())
		s.Require().NotNil(o)
		isBlacklisted := o.IsBlacklisted(s.DB)

		// assert that the organization is blacklisted
		s.True(isBlacklisted)
	})

	s.Run("should panic if the organization option returns error", func() {
		s.T().Parallel()

		errOrgOption := func(t *testing.T) fake.OrganizationOption {
			t.Helper()

			return func(o *fake.FakeOrganization) (*fake.FakeOrganization, error) {
				return o, assert.AnError
			}
		}
		newOrgFunction := func() {
			fake.NewOrganization(s.DB, errOrgOption(s.T()))
		}
		s.Panics(newOrgFunction)
	})
}
