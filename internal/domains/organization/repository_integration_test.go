package organization_test

import (
	"context"
	"database/sql"
	"time"

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
		s.ErrorContains(err, "DELETE operation on table organizations is not allowed: prevent_hard_delete_on_organizations")
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_GetOrganizationByID() {
	s.Run("should return an organization by id", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		result, err := repo.GetOrganizationByID(context.Background(), org.ID)
		s.Require().NoError(err)
		s.Equal(org.Organization, result)
	})

	s.Run("should return error when organization does not exist", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)

		_, err := repo.GetOrganizationByID(context.Background(), int64(gofakeit.Number(1000, 9999)))
		s.ErrorIs(err, sql.ErrNoRows)
	})

	s.Run("should return error when organization is deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())

		_, err := repo.GetOrganizationByID(context.Background(), org.ID)
		s.ErrorIs(err, sql.ErrNoRows)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_GetOrganizationBySubdomain() {
	s.Run("should return an organization by subdomain", func() {
		s.T().Parallel()
		r := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		result, err := r.GetOrganizationBySubdomain(context.Background(), org.Subdomain)
		s.Require().NoError(err)
		s.Equal(org.Organization, result)
	})

	s.Run("should return error when organization does not exist", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)

		_, err := repo.GetOrganizationBySubdomain(context.Background(), randomOrganizationName())
		s.ErrorIs(err, sql.ErrNoRows)
	})

	s.Run("should return error when organization is deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())

		_, err := repo.GetOrganizationBySubdomain(context.Background(), org.Subdomain)
		s.ErrorIs(err, sql.ErrNoRows)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_GetOrganizationByName() {
	s.Run("should return an organization by name", func() {
		s.T().Parallel()
		r := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		result, err := r.GetOrganizationByName(context.Background(), org.Name)
		s.Require().NoError(err)
		s.Equal(org.Organization, result)
	})

	s.Run("should return error when organization does not exist", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)

		_, err := repo.GetOrganizationByName(context.Background(), randomOrganizationName())
		s.ErrorIs(err, sql.ErrNoRows)
	})

	s.Run("should return error when organization is deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())

		_, err := repo.GetOrganizationByName(context.Background(), org.Name)
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

		result, err := repo.CreateOrganization(context.Background(), org.Subdomain, org.Name)
		s.Require().NoError(err)
		s.NotEmpty(result.ID)
		s.Equal(org.Name, result.Name)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Equal(time.UTC, result.CreatedAt.Location())
		s.Equal(time.UTC, result.UpdatedAt.Location())
		s.WithinDuration(time.Now().UTC(), result.CreatedAt, 1*time.Minute)
		s.WithinDuration(time.Now().UTC(), result.UpdatedAt, 1*time.Minute)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.DisabledAt)
		s.Nil(result.Comment)
	})

	s.Run("should return error when organization with the same subdomain already exists", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		_, err := repo.CreateOrganization(context.Background(), org.Subdomain, randomOrganizationName())

		s.Require().Error(err)
		s.ErrorContains(err, "duplicate key value violates unique constraint")
	})

	s.Run("should return error when organization with the same name already exists", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		_, err := repo.CreateOrganization(context.Background(), randomOrganizationSubdomain(), org.Name)

		s.Require().Error(err)
		s.ErrorContains(err, "duplicate key value violates unique constraint")
	})

	s.Run("should return error when organization has empty subdomain", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)

		_, err := repo.CreateOrganization(context.Background(), "", randomOrganizationName())

		s.Require().Error(err)
		s.ErrorContains(err, "violates check constraint")
	})

	s.Run("should return error when organization has empty name", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)

		_, err := repo.CreateOrganization(context.Background(), randomOrganizationSubdomain(), "")

		s.Require().Error(err)
		s.ErrorContains(err, "violates check constraint")
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_UpdateOrganization() {
	s.Run("should update an organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)
		newOrgName := randomOrganizationName()

		err := repo.UpdateOrganization(context.Background(), org.ID, newOrgName)
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Equal(newOrgName, result.Name)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.GreaterOrEqual(result.UpdatedAt, result.CreatedAt) // could be equal if the update is fast
		s.Equal(time.UTC, result.UpdatedAt.Location())
		s.WithinDuration(time.Now().UTC(), result.UpdatedAt, 1*time.Minute)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.DisabledAt)
		s.Nil(result.Comment)
	})

	s.Run("should not update an organization if already deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())
		newOrgName := randomOrganizationName()

		err := repo.UpdateOrganization(context.Background(), org.ID, newOrgName)
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Equal(org.Name, result.Name) // name should not be updated
		s.Equal(org.Subdomain, result.Subdomain)
		s.Equal(org.UpdatedAt, result.UpdatedAt) // update time should not be updated
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.NotNil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.DisabledAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_DeleteOrganization() {
	s.Run("should delete an organization", func() {
		s.T().Parallel()

		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)
		comment := gofakeit.Sentence(5)

		err := repo.DeleteOrganization(context.Background(), org.ID, comment)
		s.Require().NoError(err)

		isDeleted := org.IsDeleted(s.DB)
		s.True(isDeleted)

		result := org.FetchLatest(s.DB)
		s.Require().NotNil(result.DeletedAt)
		s.Equal(time.UTC, result.DeletedAt.Location())
		s.WithinDuration(time.Now().UTC(), *result.DeletedAt, 1*time.Minute)
		s.Require().NotNil(result.Comment)
		s.Equal(comment, *result.Comment)
	})

	s.Run("should not delete an organization if already deleted", func() {
		s.T().Parallel()

		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())
		comment := gofakeit.Sentence(5)

		err := repo.DeleteOrganization(context.Background(), org.ID, comment)
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.NotNil(result.DeletedAt)
		s.Equal(org.DeletedAt, result.DeletedAt) // delete time should not be updated
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_SuspendOrganization() {
	s.Run("should suspend an organization", func() {
		s.T().Parallel()

		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)
		comment := gofakeit.Sentence(5)

		err := repo.SuspendOrganization(context.Background(), org.ID, comment)
		s.Require().NoError(err)

		isSuspended := org.IsSuspended(s.DB)
		s.True(isSuspended)

		result := org.FetchLatest(s.DB)
		s.Require().NotNil(result.SuspendedAt)
		s.Equal(time.UTC, result.SuspendedAt.Location())
		s.WithinDuration(time.Now().UTC(), *result.SuspendedAt, 1*time.Minute)
		s.Require().NotNil(result.Comment)
		s.Equal(comment, *result.Comment)
		s.Nil(result.DeletedAt)
		s.Nil(result.DisabledAt)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.Equal(org.UpdatedAt, result.UpdatedAt)
	})

	s.Run("should not suspend an organization if already deleted", func() {
		s.T().Parallel()

		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())
		comment := gofakeit.Sentence(5)

		err := repo.SuspendOrganization(context.Background(), org.ID, comment)
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Require().Nil(result.SuspendedAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_UnsuspendOrganization() {
	s.Run("should unsuspend an organization", func() {
		s.T().Parallel()

		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationSuspended())
		comment := gofakeit.Sentence(5)

		err := repo.UnsuspendOrganization(context.Background(), org.ID, comment)
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Nil(result.SuspendedAt)
		s.Require().NotNil(result.Comment)
		s.Equal(comment, *result.Comment)
		s.Nil(result.DeletedAt)
		s.Nil(result.DisabledAt)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.Equal(org.UpdatedAt, result.UpdatedAt)
	})

	s.Run("should not unsuspend an organization if already deleted", func() {
		s.T().Parallel()

		comment := "test unsuspend comment"
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationSuspended(), fake.OrganizationDeleted())

		err := repo.UnsuspendOrganization(context.Background(), org.ID, comment)
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Require().NotNil(result.SuspendedAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_DisableOrganization() {
	s.Run("should disable an organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB)

		err := repo.DisableOrganization(context.Background(), org.ID, "test disable comment")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Require().NotNil(result.DisabledAt)
		s.Equal(time.UTC, result.DisabledAt.Location())
		s.WithinDuration(time.Now().UTC(), *result.DisabledAt, 1*time.Minute)
		s.Require().NotNil(result.Comment)
		s.Equal("test disable comment", *result.Comment)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.Equal(org.UpdatedAt, result.UpdatedAt)
	})

	s.Run("should not disable an organization if already deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDeleted())

		err := repo.DisableOrganization(context.Background(), org.ID, "test disable comment")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Require().Nil(result.DisabledAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestRepositoryIntegration_EnableOrganization() {
	s.Run("should enable an organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDisabled())

		err := repo.EnableOrganization(context.Background(), org.ID, "test enable comment")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Nil(result.DisabledAt)
		s.Require().NotNil(result.Comment)
		s.Equal("test enable comment", *result.Comment)
		s.Nil(result.DeletedAt)
		s.Nil(result.SuspendedAt)
		s.Equal(org.CreatedAt, result.CreatedAt)
		s.Equal(org.UpdatedAt, result.UpdatedAt)
	})

	s.Run("should not enable an organization if already deleted", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		org := fake.NewOrganization(s.DB, fake.OrganizationDisabled(), fake.OrganizationDeleted())

		err := repo.EnableOrganization(context.Background(), org.ID, "test enable comment")
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.Require().NotNil(result.DisabledAt)
		s.Nil(result.Comment)
	})
}
