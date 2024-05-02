package organization_test

import (
	"context"
	"database/sql"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (s *OrganizationTestSuite) TestRepositoryIntegration_GetOrganizationByID() {
	s.Run("should return error if query execution fails", func() {
		s.T().Parallel()
		db := database.NewDatabaseMock(s.T())
		r := organization.NewRepository(db)

		db.On("Get", fake.MockContext, mock.AnythingOfType("*organization.Organization"), fake.MockString, fake.MockInt64).
			Return(assert.AnError)

		_, err := r.GetOrganizationByID(context.TODO(), 1)

		s.Require().Error(err)
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should return an organization", func() {
		s.T().Parallel()
		r := organization.NewRepository(s.DB)

		org := fake.NewOrganization(s.DB)

		result, err := r.GetOrganizationByID(context.TODO(), org.ID)

		s.Require().NoError(err)
		s.Equal(org.Organization, result)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_GetOrganizationByName() {
	s.Run("should return error if query execution fails", func() {
		s.T().Parallel()
		db := database.NewDatabaseMock(s.T())
		r := organization.NewRepository(db)

		db.On("Get", fake.MockContext, mock.AnythingOfType("*organization.Organization"), fake.MockString, fake.MockString).
			Return(assert.AnError)

		_, err := r.GetOrganizationByName(context.TODO(), "TestOrg")

		s.Require().Error(err)
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should return an organization", func() {
		s.T().Parallel()
		r := organization.NewRepository(s.DB)

		org := fake.NewOrganization(s.DB)

		result, err := r.GetOrganizationByName(context.TODO(), org.Name)

		s.Require().NoError(err)
		s.Equal(org.Organization, result)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_CreateOrganization() {
	s.Run("should return error if query execution fails", func() {
		s.T().Parallel()
		db := database.NewDatabaseMock(s.T())
		r := organization.NewRepository(db)

		db.On("Exec", fake.MockContext, fake.MockInt64Ptr, fake.MockString, fake.MockString).Return(assert.AnError)

		_, err := r.CreateOrganization(context.TODO(), organization.Organization{
			Name: gofakeit.Name(),
		})

		s.Require().Error(err)
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should create an organization", func() {
		s.T().Parallel()
		r := organization.NewRepository(s.DB)
		org := organization.Organization{
			Name: gofakeit.Name(),
		}

		id, err := r.CreateOrganization(context.TODO(), org)

		s.Require().NoError(err)
		s.NotEmpty(id)

		result, err := r.GetOrganizationByID(context.TODO(), id)
		s.Require().NoError(err)
		s.Equal(org.Name, result.Name)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.BlacklistedAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_UpdateOrganization() {
	s.Run("should return error if query execution fails", func() {
		s.T().Parallel()
		db := database.NewDatabaseMock(s.T())
		r := organization.NewRepository(db)

		db.On("Exec", fake.MockContext, nil, fake.MockString, fake.MockInt64, fake.MockString).Return(assert.AnError)

		err := r.UpdateOrganization(context.TODO(), organization.Organization{
			Name: gofakeit.Name(),
		})

		s.Require().Error(err)
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should update an organization", func() {
		s.T().Parallel()
		r := organization.NewRepository(s.DB)

		org := fake.NewOrganization(s.DB)

		err := r.UpdateOrganization(context.TODO(), organization.Organization{
			ID:   org.ID,
			Name: org.Name + " Updated",
		})
		s.Require().NoError(err)

		result, err := r.GetOrganizationByID(context.TODO(), org.ID)
		s.Require().NoError(err)
		s.Equal(org.Name+" Updated", result.Name)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.GreaterOrEqual(result.UpdatedAt.Unix(), result.CreatedAt.Unix()) // could be equal if the update is fast
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.BlacklistedAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_DeleteOrganization() {
	s.Run("should return error if query execution fails", func() {
		s.T().Parallel()
		db := database.NewDatabaseMock(s.T())
		r := organization.NewRepository(db)

		db.On("Exec", fake.MockContext, nil, fake.MockString, fake.MockInt64).Return(assert.AnError)

		err := r.DeleteOrganization(context.TODO(), 1)

		s.Require().Error(err)
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should delete an organization", func() {
		s.T().Parallel()
		r := organization.NewRepository(s.DB)

		org := fake.NewOrganization(s.DB)

		err := r.DeleteOrganization(context.TODO(), org.ID)
		s.Require().NoError(err)

		_, err = r.GetOrganizationByID(context.TODO(), org.ID)
		s.Require().Error(err)
		s.ErrorIs(err, sql.ErrNoRows)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_SuspendOrganization() {
	s.Run("should return error if query execution fails", func() {
		s.T().Parallel()
		db := database.NewDatabaseMock(s.T())
		r := organization.NewRepository(db)

		db.On("Exec", fake.MockContext, nil, fake.MockString, fake.MockInt64, fake.MockString).Return(assert.AnError)

		err := r.SuspendOrganization(context.TODO(), 1, "test comment")

		s.Require().Error(err)
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should suspend an organization", func() {
		s.T().Parallel()
		r := organization.NewRepository(s.DB)

		org := fake.NewOrganization(s.DB)

		err := r.SuspendOrganization(context.TODO(), org.ID, "test comment")
		s.Require().NoError(err)

		result, err := r.GetOrganizationByID(context.TODO(), org.ID)
		s.Require().NoError(err)
		s.NotNil(result.SuspendedAt)
		s.Require().NotNil(result.Comment)
		s.Equal("test comment", *result.Comment)
		s.Nil(result.DeletedAt)
		s.Nil(result.BlacklistedAt)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.Equal(org.UpdatedAt, result.UpdatedAt)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_UnsuspendOrganization() {
	s.Run("should return error if query execution fails", func() {
		s.T().Parallel()
		db := database.NewDatabaseMock(s.T())
		r := organization.NewRepository(db)

		db.On("Exec", fake.MockContext, nil, fake.MockString, fake.MockInt64, fake.MockString).Return(assert.AnError)

		err := r.UnsuspendOrganization(context.TODO(), 1, "test comment")

		s.Require().Error(err)
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should unsuspend an organization", func() {
		s.T().Parallel()
		r := organization.NewRepository(s.DB)

		org := fake.NewOrganization(s.DB, fake.OrganizationSuspended())

		err := r.UnsuspendOrganization(context.TODO(), org.ID, "test comment")
		s.Require().NoError(err)

		result, err := r.GetOrganizationByID(context.TODO(), org.ID)
		s.Require().NoError(err)
		s.Nil(result.SuspendedAt)
		s.Require().NotNil(result.Comment)
		s.Equal("test comment", *result.Comment)
		s.Nil(result.DeletedAt)
		s.Nil(result.BlacklistedAt)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.Equal(org.UpdatedAt, result.UpdatedAt)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_BlacklistOrganization() {
	s.Run("should return error if query execution fails", func() {
		s.T().Parallel()
		db := database.NewDatabaseMock(s.T())
		r := organization.NewRepository(db)

		db.On("Exec", fake.MockContext, nil, fake.MockString, fake.MockInt64, fake.MockString).Return(assert.AnError)

		err := r.BlacklistOrganization(context.TODO(), 1, "test comment")

		s.Require().Error(err)
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should blacklist an organization", func() {
		s.T().Parallel()
		r := organization.NewRepository(s.DB)

		org := fake.NewOrganization(s.DB)

		err := r.BlacklistOrganization(context.TODO(), org.ID, "test comment")
		s.Require().NoError(err)

		result, err := r.GetOrganizationByID(context.TODO(), org.ID)
		s.Require().NoError(err)
		s.NotNil(result.BlacklistedAt)
		s.Require().NotNil(result.Comment)
		s.Equal("test comment", *result.Comment)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.Equal(org.UpdatedAt, result.UpdatedAt)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_UnblacklistOrganization() {
	s.Run("should return error if query execution fails", func() {
		s.T().Parallel()
		db := database.NewDatabaseMock(s.T())
		r := organization.NewRepository(db)

		db.On("Exec", fake.MockContext, nil, fake.MockString, fake.MockInt64, fake.MockString).Return(assert.AnError)

		err := r.UnblacklistOrganization(context.TODO(), 1, "test comment")

		s.Require().Error(err)
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should unblacklist an organization", func() {
		s.T().Parallel()
		r := organization.NewRepository(s.DB)

		org := fake.NewOrganization(s.DB, fake.OrganizationBlacklisted())

		err := r.UnblacklistOrganization(context.TODO(), org.ID, "test comment")
		s.Require().NoError(err)

		result, err := r.GetOrganizationByID(context.TODO(), org.ID)
		s.Require().NoError(err)
		s.Nil(result.BlacklistedAt)
		s.Require().NotNil(result.Comment)
		s.Equal("test comment", *result.Comment)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.Equal(org.UpdatedAt, result.UpdatedAt)
	})
}
