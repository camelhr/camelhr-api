package organization_test

import (
	"context"
	"database/sql"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
)

func (s *OrganizationTestSuite) TestServiceIntegration_GetOrganizationByID() {
	s.Run("should return organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := organization.Organization{
			Name: gofakeit.Name(),
		}

		id, err := svc.CreateOrganization(context.TODO(), org)
		s.Require().NoError(err)

		result, err := svc.GetOrganizationByID(context.TODO(), id)
		s.Require().NoError(err)
		s.Equal(id, result.ID)
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

func (s *OrganizationTestSuite) TestServiceIntegration_GetOrganizationByName() {
	s.Run("should return organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := organization.Organization{
			Name: gofakeit.Name(),
		}

		id, err := svc.CreateOrganization(context.TODO(), org)
		s.Require().NoError(err)
		s.NotEmpty(id)

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
	s.Run("should return organization ID", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := organization.Organization{
			Name: gofakeit.Name(),
		}

		id, err := svc.CreateOrganization(context.TODO(), org)

		s.Require().NoError(err)
		s.NotEmpty(id)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_UpdateOrganization() {
	s.Run("should update organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := organization.Organization{
			Name: gofakeit.Name(),
		}

		id, err := svc.CreateOrganization(context.TODO(), org)
		s.Require().NoError(err)

		updateOrg := organization.Organization{
			ID:   id,
			Name: "UpdatedOrg",
		}
		err = svc.UpdateOrganization(context.TODO(), updateOrg)
		s.Require().NoError(err)

		result, err := svc.GetOrganizationByID(context.TODO(), id)
		s.Require().NoError(err)
		s.Equal(updateOrg.Name, result.Name)
		s.Nil(result.DeletedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.GreaterOrEqual(result.UpdatedAt.Unix(), result.CreatedAt.Unix()) // could be equal if the update is fast
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
		org := organization.Organization{
			Name: gofakeit.Name(),
		}

		id, err := svc.CreateOrganization(context.TODO(), org)
		s.Require().NoError(err)

		err = svc.DeleteOrganization(context.TODO(), id)
		s.Require().NoError(err)

		_, err = svc.GetOrganizationByID(context.TODO(), id)
		s.ErrorIs(err, sql.ErrNoRows)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_SuspendOrganization() {
	s.Run("should suspend organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo)
		org := organization.Organization{
			Name: gofakeit.Name(),
		}

		id, err := svc.CreateOrganization(context.TODO(), org)
		s.Require().NoError(err)

		err = svc.SuspendOrganization(context.TODO(), id, "test suspend")
		s.Require().NoError(err)

		result, err := svc.GetOrganizationByID(context.TODO(), id)
		s.Require().NoError(err)
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
		org := organization.Organization{
			Name: gofakeit.Name(),
		}

		id, err := svc.CreateOrganization(context.TODO(), org)
		s.Require().NoError(err)

		err = svc.SuspendOrganization(context.TODO(), id, "test suspend")
		s.Require().NoError(err)

		err = svc.UnsuspendOrganization(context.TODO(), id, "test unsuspend")
		s.Require().NoError(err)

		result, err := svc.GetOrganizationByID(context.TODO(), id)
		s.Require().NoError(err)
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
		org := organization.Organization{
			Name: gofakeit.Name(),
		}

		id, err := svc.CreateOrganization(context.TODO(), org)
		s.Require().NoError(err)

		err = svc.BlacklistOrganization(context.TODO(), id, "test blacklist")
		s.Require().NoError(err)

		result, err := svc.GetOrganizationByID(context.TODO(), id)
		s.Require().NoError(err)
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
		org := organization.Organization{
			Name: gofakeit.Name(),
		}

		id, err := svc.CreateOrganization(context.TODO(), org)
		s.Require().NoError(err)

		err = svc.BlacklistOrganization(context.TODO(), id, "test blacklist")
		s.Require().NoError(err)

		err = svc.UnblacklistOrganization(context.TODO(), id, "test unblacklist")
		s.Require().NoError(err)

		result, err := svc.GetOrganizationByID(context.TODO(), id)
		s.Require().NoError(err)
		s.Require().NotNil(result.Comment)
		s.Equal("test unblacklist", *result.Comment)
		s.Nil(result.BlacklistedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
	})
}
