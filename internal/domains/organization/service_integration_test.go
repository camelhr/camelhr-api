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

		result, err := svc.GetOrganizationByID(context.Background(), org.ID)
		s.Require().NoError(err)
		s.Equal(org.ID, result.ID)
		s.Equal(org.Name, result.Name)
		s.Nil(result.DeletedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.DisabledAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_GetOrganizationBySubdomain() {
	s.Run("should return organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB)

		result, err := svc.GetOrganizationBySubdomain(context.Background(), org.Subdomain)
		s.Require().NoError(err)
		s.Equal(org.Subdomain, result.Subdomain)
		s.Nil(result.DeletedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.DisabledAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_GetOrganizationByName() {
	s.Run("should return organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB)

		result, err := svc.GetOrganizationByName(context.Background(), org.Name)
		s.Require().NoError(err)
		s.Equal(org.Name, result.Name)
		s.Nil(result.DeletedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.DisabledAt)
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

		result, err := svc.CreateOrganization(context.Background(), org.Subdomain, org.Name)
		s.Require().NoError(err)
		s.NotZero(result.ID)
		s.Equal(org.Name, result.Name)
		s.Nil(result.DeletedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.DisabledAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_UpdateOrganization() {
	s.Run("should update organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB)
		newOrgName := randomOrganizationName()

		err := svc.UpdateOrganization(context.Background(), org.ID, newOrgName)
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Equal(newOrgName, result.Name)
		s.Nil(result.DeletedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.GreaterOrEqual(result.UpdatedAt, result.CreatedAt) // could be equal if the update is fast
		s.Nil(result.SuspendedAt)
		s.Nil(result.DisabledAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_DeleteOrganization() {
	s.Run("should delete organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB)

		err := svc.DeleteOrganization(context.Background(), org.ID)
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

		err := svc.SuspendOrganization(context.Background(), org.ID, "test suspend")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Require().NotNil(result.Comment)
		s.Equal("test suspend", *result.Comment)
		s.NotNil(result.SuspendedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Nil(result.DeletedAt)
		s.Nil(result.DisabledAt)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_UnsuspendOrganization() {
	s.Run("should unsuspend organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB, fake.OrganizationSuspended())

		err := svc.UnsuspendOrganization(context.Background(), org.ID, "test unsuspend")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Require().NotNil(result.Comment)
		s.Equal("test unsuspend", *result.Comment)
		s.Nil(result.SuspendedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Nil(result.DeletedAt)
		s.Nil(result.DisabledAt)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_DisableOrganization() {
	s.Run("should disable organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB)

		err := svc.DisableOrganization(context.Background(), org.ID, "test disable")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Require().NotNil(result.Comment)
		s.Equal("test disable", *result.Comment)
		s.NotNil(result.DisabledAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_EnableOrganization() {
	s.Run("should enable organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := fake.NewOrganization(s.DB, fake.OrganizationDisabled())

		err := svc.EnableOrganization(context.Background(), org.ID, "test enable")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Require().NotNil(result.Comment)
		s.Equal("test enable", *result.Comment)
		s.Nil(result.DisabledAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
	})
}
