package user_test

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
)

func (s *UserTestSuite) TestServiceIntegration_GetUserByID() {
	s.Run("should return user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		result, err := svc.GetUserByID(context.Background(), u.ID)
		s.Require().NoError(err)
		s.Equal(u.ID, result.ID)
		s.Equal(u.Email, result.Email)
		s.Nil(result.DisabledAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.Comment)
	})
}

func (s *UserTestSuite) TestServiceIntegration_GetUserByAPIToken() {
	s.Run("should return user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)
		s.Require().NotNil(u.APIToken)

		result, err := svc.GetUserByAPIToken(context.Background(), *u.APIToken)
		s.Require().NoError(err)
		s.Equal(u.ID, result.ID)
		s.Equal(u.Email, result.Email)
		s.Nil(result.DisabledAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.Comment)
	})
}

func (s *UserTestSuite) TestServiceIntegration_GetUserByOrgIDEmail() {
	s.Run("should return user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		result, err := svc.GetUserByOrgIDEmail(context.Background(), o.ID, u.Email)
		s.Require().NoError(err)
		s.Equal(u.ID, result.ID)
		s.Equal(u.Email, result.Email)
		s.Nil(result.DisabledAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.Comment)
	})
}

func (s *UserTestSuite) TestServiceIntegration_GetUserByOrgSubdomainEmail() {
	s.Run("should return user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		result, err := svc.GetUserByOrgSubdomainEmail(context.Background(), o.Subdomain, u.Email)
		s.Require().NoError(err)
		s.Equal(u.ID, result.ID)
		s.Equal(u.Email, result.Email)
		s.Nil(result.DisabledAt)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.Comment)
	})
}

func (s *UserTestSuite) TestServiceIntegration_CreateUser() {
	s.Run("should create user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo)
		o := fake.NewOrganization(s.DB)
		email := gofakeit.Email()
		password := gofakeit.Password(true, true, true, true, false, 12)

		result, err := svc.CreateUser(context.Background(), o.ID, email, password)
		s.Require().NoError(err)
		s.Equal(email, result.Email)
		s.Nil(result.DisabledAt)
		s.Nil(result.Comment)
		s.Nil(result.APIToken)
		s.False(result.IsOwner)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.DeletedAt)
	})
}

func (s *UserTestSuite) TestServiceIntegration_CreateOwner() {
	s.Run("should create owner", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo)
		o := fake.NewOrganization(s.DB)
		email := gofakeit.Email()
		password := gofakeit.Password(true, true, true, true, false, 12)

		result, err := svc.CreateOwner(context.Background(), o.ID, email, password)
		s.Require().NoError(err)
		s.Equal(email, result.Email)
		s.Nil(result.DisabledAt)
		s.Nil(result.Comment)
		s.Nil(result.APIToken)
		s.True(result.IsOwner)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Nil(result.DeletedAt)
	})
}

func (s *UserTestSuite) TestServiceIntegration_ResetPassword() {
	s.Run("should reset password", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)
		newPassword := gofakeit.Password(true, true, true, true, false, 12)

		err := svc.ResetPassword(context.Background(), u.ID, newPassword)
		s.Require().NoError(err)

		result, err := svc.GetUserByID(context.Background(), u.ID)
		s.Require().NoError(err)
		s.NotEqual(u.PasswordHash, result.PasswordHash)
	})
}

func (s *UserTestSuite) TestServiceIntegration_DeleteUser() {
	s.Run("should delete user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		err := svc.DeleteUser(context.Background(), u.ID)
		s.Require().NoError(err)

		result := u.IsDeleted(s.DB)
		s.True(result)
	})
}

func (s *UserTestSuite) TestServiceIntegration_DisableUser() {
	s.Run("should disable user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		err := svc.DisableUser(context.Background(), u.ID, "test")
		s.Require().NoError(err)

		result, err := svc.GetUserByID(context.Background(), u.ID)
		s.Require().NoError(err)
		s.NotNil(result.DisabledAt)
		s.Require().NotNil(result.Comment)
		s.Equal("test", *result.Comment)
	})
}

func (s *UserTestSuite) TestServiceIntegration_EnableUser() {
	s.Run("should enable user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDisabled())

		err := svc.EnableUser(context.Background(), u.ID, "test")
		s.Require().NoError(err)

		result, err := svc.GetUserByID(context.Background(), u.ID)
		s.Require().NoError(err)
		s.Nil(result.DisabledAt)
		s.Require().NotNil(result.Comment)
		s.Equal("test", *result.Comment)
	})
}

func (s *UserTestSuite) TestServiceIntegration_GenerateAPIToken() {
	s.Run("should generate API token", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		err := svc.GenerateAPIToken(context.Background(), u.ID)
		s.Require().NoError(err)

		result, err := svc.GetUserByID(context.Background(), u.ID)
		s.Require().NoError(err)
		s.NotNil(result.APIToken)
	})
}

func (s *UserTestSuite) TestServiceIntegration_ResetAPIToken() {
	s.Run("should reset API token", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		err := svc.ResetAPIToken(context.Background(), u.ID)
		s.Require().NoError(err)

		result, err := svc.GetUserByID(context.Background(), u.ID)
		s.Require().NoError(err)
		s.NotNil(result.APIToken)
	})
}

func (s *UserTestSuite) TestServiceIntegration_SetEmailVerified() {
	s.Run("should set email verified", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserEmailNotVerified())

		err := svc.SetEmailVerified(context.Background(), u.ID)
		s.Require().NoError(err)

		result, err := svc.GetUserByID(context.Background(), u.ID)
		s.Require().NoError(err)
		s.True(result.IsEmailVerified)
	})
}
