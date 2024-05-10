package fake_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
	"github.com/stretchr/testify/assert"
)

func (s *FakeTestSuite) TestFakeUser() {
	s.Run("should create a user with default values", func() {
		s.T().Parallel()

		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		// assert that the user is created with the default values
		s.Require().NotNil(u)
		s.NotEmpty(u.ID)
		s.Equal(o.ID, u.OrganizationID)
		s.NotEmpty(u.Email)
		s.NotEmpty(u.PasswordHash)
		s.NotNil(u.APIToken)
		s.False(u.IsOwner)
		s.Nil(u.DisabledAt)
		s.Nil(u.Comment)
		s.NotZero(u.CreatedAt)
		s.NotZero(u.UpdatedAt)
		s.Equal(u.CreatedAt, u.UpdatedAt)
		s.Equal(time.UTC, u.CreatedAt.Location())
		s.Equal(time.UTC, u.UpdatedAt.Location())
		s.WithinDuration(time.Now().UTC(), u.CreatedAt, 1*time.Minute)
		s.WithinDuration(time.Now().UTC(), u.UpdatedAt, 1*time.Minute)
		s.Nil(u.DeletedAt)
	})

	s.Run("should create a user without an API token", func() {
		s.T().Parallel()

		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserWithoutToken())

		// assert that the user is created without an API token
		s.Require().NotNil(u)
		s.Nil(u.APIToken)
	})

	s.Run("should create a user with custom email", func() {
		s.T().Parallel()

		o := fake.NewOrganization(s.DB)
		email := gofakeit.Email()
		u := fake.NewUser(s.DB, o.ID, fake.UserEmail(email))

		// assert that the user is created with the specified email
		s.Require().NotNil(u)
		s.Equal(email, u.Email)
	})

	s.Run("should create a user as the owner of the organization", func() {
		s.T().Parallel()

		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserIsOwner())

		// assert that the user is created as the owner of the organization
		s.Require().NotNil(u)
		s.True(u.IsOwner)
	})

	s.Run("should create a disabled user", func() {
		s.T().Parallel()

		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDisabled())

		// assert that the user is created as disabled
		s.Require().NotNil(u)
		s.NotNil(u.DisabledAt)
		s.Equal(time.UTC, u.DisabledAt.Location())
		s.WithinDuration(time.Now().UTC(), *u.DisabledAt, 1*time.Minute)
	})

	s.Run("should create a deleted user", func() {
		s.T().Parallel()

		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDeleted())

		// assert that the user is set deleted
		s.Require().NotNil(u)
		s.Require().NotNil(u.DeletedAt)
		s.Equal(time.UTC, u.DeletedAt.Location())
		s.WithinDuration(time.Now().UTC(), *u.DeletedAt, 1*time.Minute)
	})

	s.Run("should return true if user is deleted", func() {
		s.T().Parallel()

		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDeleted())
		s.Require().NotNil(u)

		isDeleted := u.IsDeleted(s.DB)

		// assert that the user is deleted
		s.True(isDeleted)
	})

	s.Run("should return true if user is disabled", func() {
		s.T().Parallel()

		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDisabled())
		s.Require().NotNil(u)

		isDisabled := u.IsDisabled(s.DB)

		// assert that the user is disabled
		s.True(isDisabled)
	})

	s.Run("should panic if an error occurs while creating a user", func() {
		s.T().Parallel()

		errUserOption := func(t *testing.T) fake.UserOption {
			t.Helper()

			return func(o *fake.FakeUser) (*fake.FakeUser, error) {
				return o, assert.AnError
			}
		}

		newUserFunc := func() {
			fake.NewUser(s.DB, gofakeit.Int64(), errUserOption(s.T()))
		}
		s.Panics(newUserFunc)
	})
}
