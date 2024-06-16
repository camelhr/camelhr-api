package organization_test

import (
	"context"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/domains/session"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
)

func (s *OrganizationTestSuite) TestServiceIntegration_GetOrganizationByID() {
	s.Run("should return organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo, nil)
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
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_GetOrganizationBySubdomain() {
	s.Run("should return organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo, nil)
		org := fake.NewOrganization(s.DB)

		result, err := svc.GetOrganizationBySubdomain(context.Background(), org.Subdomain)
		s.Require().NoError(err)
		s.Equal(org.Subdomain, result.Subdomain)
		s.Nil(result.DeletedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_GetOrganizationByName() {
	s.Run("should return organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo, nil)
		org := fake.NewOrganization(s.DB)

		result, err := svc.GetOrganizationByName(context.Background(), org.Name)
		s.Require().NoError(err)
		s.Equal(org.Name, result.Name)
		s.Nil(result.DeletedAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.SuspendedAt)
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_CreateOrganization() {
	s.Run("should create organization with default values", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo, nil)
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
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_UpdateOrganization() {
	s.Run("should update organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo, nil)
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
		s.Nil(result.Comment)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_DeleteOrganization() {
	s.Run("should delete organization", func() {
		s.T().Parallel()

		sessionManager := session.NewRedisSessionManager(s.RedisClient)
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo, sessionManager)
		org := fake.NewOrganization(s.DB)
		u1 := fake.NewUser(s.DB, org.ID)
		u2 := fake.NewUser(s.DB, org.ID)
		comment := gofakeit.Sentence(5)

		// create session for users
		sessionKey := fmt.Sprintf("session:org:%v:user:%v", org.ID, u1.ID)
		err := s.RedisClient.HSet(context.Background(), sessionKey, "jwt", gofakeit.UUID()).Err()
		s.Require().NoError(err)

		sessionKey = fmt.Sprintf("session:org:%v:user:%v", org.ID, u2.ID)
		err = s.RedisClient.HSet(context.Background(), sessionKey, "jwt", gofakeit.UUID()).Err()
		s.Require().NoError(err)

		err = svc.DeleteOrganization(context.Background(), org.ID, comment)
		s.Require().NoError(err)

		result := org.FetchLatest(s.DB)
		s.NotNil(result.DeletedAt)

		// check if all user sessions are deleted
		sessionKey = fmt.Sprintf("session:org:%v:user:%v", org.ID, u1.ID)
		exist := s.RedisClient.Exists(context.Background(), sessionKey).Val()
		s.Zero(exist)

		sessionKey = fmt.Sprintf("session:org:%v:user:%v", org.ID, u2.ID)
		exist = s.RedisClient.Exists(context.Background(), sessionKey).Val()
		s.Zero(exist)
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_SuspendOrganization() {
	s.Run("should suspend organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo, nil)
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
	})
}

func (s *OrganizationTestSuite) TestServiceIntegration_UnsuspendOrganization() {
	s.Run("should unsuspend organization", func() {
		s.T().Parallel()
		repo := organization.NewRepository(s.DB)
		svc := organization.NewService(repo, nil)
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
	})
}
