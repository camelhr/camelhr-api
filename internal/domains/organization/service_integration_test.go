package organization_test

import (
	"context"

	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
)

func (s *OrganizationTestSuite) TestServiceIntegration_GetOrganizationByID() {
	s.Run("should return organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB)

		result, err := svc.GetOrganizationByID(context.TODO(), org.ID)
		s.Require().NoError(err)
		s.Equal(org.ID, result.ID)
		s.Equal(org.Name, result.Name)
		s.Nil(result.DeletedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.BlacklistedAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_GetOrganizationBySubdomain() {
	s.Run("should return organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB)

		result, err := svc.GetOrganizationBySubdomain(context.TODO(), org.Subdomain)
		s.Require().NoError(err)
		s.Equal(org.Subdomain, result.Subdomain)
		s.Nil(result.DeletedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.BlacklistedAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_GetOrganizationByName() {
	s.Run("should return organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB)

		result, err := svc.GetOrganizationByName(context.TODO(), org.Name)
		s.Require().NoError(err)
		s.Equal(org.Name, result.Name)
		s.Nil(result.DeletedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.BlacklistedAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_CreateOrganization() {
	s.Run("should create organization with default values", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := organization.Organization{
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		id, err := svc.CreateOrganization(context.TODO(), org)
		s.Require().NoError(err)

		result, err := svc.GetOrganizationByID(context.TODO(), id)
		s.Require().NoError(err)
		s.Equal(org.Name, result.Name)
		s.Nil(result.DeletedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.BlacklistedAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_UpdateOrganization() {
	s.Run("should update organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB)

		updateOrg := organization.Organization{
			ID:        org.ID,
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}
		err := svc.UpdateOrganization(context.TODO(), updateOrg)
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Equal(updateOrg.Name, result.Name)
		s.Nil(result.DeletedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.GreaterOrEqual(result.UpdatedAt, result.CreatedAt) // could be equal if the update is fast
		s.Nil(result.SuspendedAt)
		s.Nil(result.BlacklistedAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_DeleteOrganization() {
	s.Run("should delete organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB)

		err := svc.DeleteOrganization(context.TODO(), org.ID)
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.NotNil(result.DeletedAt)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_SuspendOrganization() {
	s.Run("should suspend organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB)

		err := svc.SuspendOrganization(context.TODO(), org.ID, "test suspend")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Require().NotNil(result.Comment)
		s.Equal("test suspend", *result.Comment)
		s.NotNil(result.SuspendedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Nil(result.DeletedAt)
		s.Nil(result.BlacklistedAt)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_UnsuspendOrganization() {
	s.Run("should unsuspend organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB, fake.OrganizationSuspended())

		err := svc.UnsuspendOrganization(context.TODO(), org.ID, "test unsuspend")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Require().NotNil(result.Comment)
		s.Equal("test unsuspend", *result.Comment)
		s.Nil(result.SuspendedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Nil(result.DeletedAt)
		s.Nil(result.BlacklistedAt)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_BlacklistOrganization() {
	s.Run("should blacklist organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB)

		err := svc.BlacklistOrganization(context.TODO(), org.ID, "test blacklist")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Require().NotNil(result.Comment)
		s.Equal("test blacklist", *result.Comment)
		s.NotNil(result.BlacklistedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_UnblacklistOrganization() {
	s.Run("should unblacklist organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB, fake.OrganizationBlacklisted())

		err := svc.UnblacklistOrganization(context.TODO(), org.ID, "test unblacklist")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Require().NotNil(result.Comment)
		s.Equal("test unblacklist", *result.Comment)
		s.Nil(result.BlacklistedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
	})
}
