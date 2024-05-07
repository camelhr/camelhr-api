package organization_test

import (
	"context"
	"database/sql"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
)

func (s *OrganizationTestSuite) TestRepositoryIntegration_ForbiddenOperations() {
	s.Run("should return error when organization is truncated", func() {
		s.T().Parallel()
		err := s.DB.Exec(context.Background(), nil, "TRUNCATE organizations CASCADE")
		s.Require().Error(err)
		s.ErrorContains(err, "TRUNCATE operation on table organizations is not allowed: prevent_truncate_on_organizations")
	})

	s.Run("should return error when delete query is performed on organizations", func() {
		s.T().Parallel()
		fake.NewOrganization(s.DB)
		err := s.DB.Exec(context.Background(), nil, "DELETE FROM organizations WHERE organization_id = 1")
		s.Require().Error(err)
		s.ErrorContains(err, "DELETE operation on table organizations is not allowed: prevent_delete_on_organizations")
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_GetOrganizationByID() {
	s.Run("should return an organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		result, err := repo.GetOrganizationByID(context.TODO(), org.ID)
		s.Require().NoError(err)
		s.Equal(org.Organization, result)
	})

	s.Run("should return error when organization does not exist", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)

		_, err := repo.GetOrganizationByID(context.TODO(), int64(gofakeit.Number(1000, 9999)))
		s.ErrorIs(err, sql.ErrNoRows)
	})

	s.Run("should return error when organization is deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())

		_, err := repo.GetOrganizationByID(context.TODO(), org.ID)
		s.ErrorIs(err, sql.ErrNoRows)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_GetOrganizationBySubdomain() {
	s.Run("should return an organization", func() {
		s.T().Parallel()
		r := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		result, err := r.GetOrganizationBySubdomain(context.TODO(), org.Subdomain)
		s.Require().NoError(err)
		s.Equal(org.Organization, result)
	})

	s.Run("should return error when organization does not exist", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)

		_, err := repo.GetOrganizationBySubdomain(context.TODO(), randomOrganizationName())
		s.ErrorIs(err, sql.ErrNoRows)
	})

	s.Run("should return error when organization is deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())

		_, err := repo.GetOrganizationBySubdomain(context.TODO(), org.Subdomain)
		s.ErrorIs(err, sql.ErrNoRows)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_GetOrganizationByName() {
	s.Run("should return an organization", func() {
		s.T().Parallel()
		r := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		result, err := r.GetOrganizationByName(context.TODO(), org.Name)
		s.Require().NoError(err)
		s.Equal(org.Organization, result)
	})

	s.Run("should return error when organization does not exist", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)

		_, err := repo.GetOrganizationByName(context.TODO(), randomOrganizationName())
		s.ErrorIs(err, sql.ErrNoRows)
	})

	s.Run("should return error when organization is deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())

		_, err := repo.GetOrganizationByName(context.TODO(), org.Name)
		s.ErrorIs(err, sql.ErrNoRows)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_CreateOrganization() {
	s.Run("should create an organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := organization.Organization{
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		id, err := repo.CreateOrganization(context.TODO(), org)
		s.Require().NoError(err)
		s.NotEmpty(id)

		result, err := repo.GetOrganizationByID(context.TODO(), id)
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

	s.Run("should return error when organization with the same subdomain already exists", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		_, err := repo.CreateOrganization(context.TODO(), organization.Organization{
			Subdomain: org.Subdomain,
			Name:      randomOrganizationName(),
		})

		s.Require().Error(err)
		s.ErrorContains(err, "duplicate key value violates unique constraint")
	})

	s.Run("should return error when organization with the same name already exists", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		_, err := repo.CreateOrganization(context.TODO(), organization.Organization{
			Subdomain: randomOrganizationSubdomain(),
			Name:      org.Name,
		})

		s.Require().Error(err)
		s.ErrorContains(err, "duplicate key value violates unique constraint")
	})

	s.Run("should return error when organization has empty subdomain", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)

		_, err := repo.CreateOrganization(context.TODO(), organization.Organization{
			Subdomain: "",
			Name:      randomOrganizationName(),
		})

		s.Require().Error(err)
		s.ErrorContains(err, "violates check constraint")
	})

	s.Run("should return error when organization has empty name", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)

		_, err := repo.CreateOrganization(context.TODO(), organization.Organization{
			Subdomain: randomOrganizationSubdomain(),
			Name:      "",
		})

		s.Require().Error(err)
		s.ErrorContains(err, "violates check constraint")
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_UpdateOrganization() {
	s.Run("should update an organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		err := repo.UpdateOrganization(context.TODO(), organization.Organization{
			ID:        org.ID,
			Subdomain: "updated-subdomain",
			Name:      "updated name",
		})
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Equal("updated-subdomain", result.Subdomain)
		s.Equal("updated name", result.Name)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.GreaterOrEqual(result.UpdatedAt, result.CreatedAt) // could be equal if the update is fast
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.BlacklistedAt)
		s.Nil(result.Comment)
	})

	s.Run("should not update an organization if already deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())

		err := repo.UpdateOrganization(context.TODO(), organization.Organization{
			ID:        org.ID,
			Subdomain: "delete-update-subdomain",
			Name:      "delete update org",
		})
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Equal(org.Subdomain, result.Subdomain) // subdomain should not be updated
		s.Equal(org.Name, result.Name)           // name should not be updated
		s.Equal(org.UpdatedAt, result.UpdatedAt) // update time should not be updated
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.NotNil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.BlacklistedAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_DeleteOrganization() {
	s.Run("should delete an organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		err := repo.DeleteOrganization(context.TODO(), org.ID)
		s.Require().NoError(err)

		isDeleted := org.IsDeleted(s.DB)
		s.True(isDeleted)
	})

	s.Run("should not delete an organization if already deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())

		err := repo.DeleteOrganization(context.TODO(), org.ID)
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.NotNil(result.DeletedAt)
		s.Equal(org.DeletedAt, result.DeletedAt) // delete time should not be updated
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_SuspendOrganization() {
	s.Run("should suspend an organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		err := repo.SuspendOrganization(context.TODO(), org.ID, "test suspend comment")
		s.Require().NoError(err)

		isSuspended := org.IsSuspended(s.DB)
		s.True(isSuspended)

		result := org.FetchLatest(s.DB)
		s.NotNil(result.SuspendedAt)
		s.Require().NotNil(result.Comment)
		s.Equal("test suspend comment", *result.Comment)
		s.Nil(result.DeletedAt)
		s.Nil(result.BlacklistedAt)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.Equal(org.UpdatedAt, result.UpdatedAt)
	})

	s.Run("should not suspend an organization if already deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())

		err := repo.SuspendOrganization(context.TODO(), org.ID, "test suspend comment")
		s.Require().NoError(err)

		isSuspended := org.IsSuspended(s.DB)
		s.False(isSuspended)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_UnsuspendOrganization() {
	s.Run("should unsuspend an organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationSuspended())

		err := repo.UnsuspendOrganization(context.TODO(), org.ID, "test unsuspend comment")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Nil(result.SuspendedAt)
		s.Require().NotNil(result.Comment)
		s.Equal("test unsuspend comment", *result.Comment)
		s.Nil(result.DeletedAt)
		s.Nil(result.BlacklistedAt)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.Equal(org.UpdatedAt, result.UpdatedAt)
	})

	s.Run("should not unsuspend an organization if already deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationSuspended(), fake.OrganizationDeleted())

		err := repo.UnsuspendOrganization(context.TODO(), org.ID, "test unsuspend comment")
		s.Require().NoError(err)

		isSuspended := org.IsSuspended(s.DB)
		s.True(isSuspended)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_BlacklistOrganization() {
	s.Run("should blacklist an organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		err := repo.BlacklistOrganization(context.TODO(), org.ID, "test blacklist comment")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.NotNil(result.BlacklistedAt)
		s.Require().NotNil(result.Comment)
		s.Equal("test blacklist comment", *result.Comment)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.Equal(org.UpdatedAt, result.UpdatedAt)
	})

	s.Run("should not blacklist an organization if already deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())

		err := repo.BlacklistOrganization(context.TODO(), org.ID, "test blacklist comment")
		s.Require().NoError(err)

		isBlacklisted := org.IsBlacklisted(s.DB)
		s.False(isBlacklisted)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_UnblacklistOrganization() {
	s.Run("should unblacklist an organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationBlacklisted())

		err := repo.UnblacklistOrganization(context.TODO(), org.ID, "test unblacklist comment")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Nil(result.BlacklistedAt)
		s.Require().NotNil(result.Comment)
		s.Equal("test unblacklist comment", *result.Comment)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.Equal(org.UpdatedAt, result.UpdatedAt)
	})

	s.Run("should not unblacklist an organization if already deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationBlacklisted(), fake.OrganizationDeleted())

		err := repo.UnblacklistOrganization(context.TODO(), org.ID, "test unblacklist comment")
		s.Require().NoError(err)

		isBlacklisted := org.IsBlacklisted(s.DB)
		s.True(isBlacklisted)
	})
}
