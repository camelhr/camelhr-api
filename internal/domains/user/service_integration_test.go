package user_test

import (
	"context"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/session"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
)

func (s *UserTestSuite) TestServiceIntegration_GetUserByID() {
	s.Run("should return user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo, nil)
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
		svc := user.NewService(repo, nil)
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

func (s *UserTestSuite) TestServiceIntegration_GetUserByOrgSubdomainAPIToken() {
	s.Run("should return user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo, nil)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)
		s.Require().NotNil(u.APIToken)

		result, err := svc.GetUserByOrgSubdomainAPIToken(context.Background(), o.Subdomain, *u.APIToken)
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
		svc := user.NewService(repo, nil)
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
		svc := user.NewService(repo, nil)
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
		svc := user.NewService(repo, nil)
		o := fake.NewOrganization(s.DB)
		email := gofakeit.Email()
		password := generatePassword()

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
		svc := user.NewService(repo, nil)
		o := fake.NewOrganization(s.DB)
		email := gofakeit.Email()
		password := generatePassword()

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
		svc := user.NewService(repo, nil)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)
		newPassword := generatePassword()

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
		sessionManager := session.NewRedisSessionManager(s.RedisClient)
		svc := user.NewService(repo, sessionManager)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)
		sessionKey := fmt.Sprintf("session:org:%v:user:%v", o.ID, u.ID)
		s.RedisClient.HSet(context.Background(), sessionKey, "jwt", gofakeit.UUID())

		err := svc.DeleteUser(context.Background(), u.ID)
		s.Require().NoError(err)

		result := u.IsDeleted(s.DB)
		s.True(result)

		// check if the user's session is deleted
		exist := s.RedisClient.Exists(context.Background(), sessionKey).Val()
		s.Zero(exist)
	})
}

func (s *UserTestSuite) TestServiceIntegration_DeleteAllUsersByOrgID() {
	s.Run("should delete all users by org id", func() {
		s.T().Parallel()

		sessionManager := session.NewRedisSessionManager(s.RedisClient)
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo, sessionManager)
		o := fake.NewOrganization(s.DB)
		u1 := fake.NewUser(s.DB, o.ID)
		sessionKey1 := fmt.Sprintf("session:org:%v:user:%v", o.ID, u1.ID)
		s.RedisClient.HSet(context.Background(), sessionKey1, "jwt", gofakeit.UUID())
		u2 := fake.NewUser(s.DB, o.ID)
		sessionKey2 := fmt.Sprintf("session:org:%v:user:%v", o.ID, u2.ID)
		s.RedisClient.HSet(context.Background(), sessionKey2, "jwt", gofakeit.UUID())

		err := svc.DeleteAllUsersByOrgID(context.Background(), o.ID)
		s.Require().NoError(err)

		result1 := u1.IsDeleted(s.DB)
		s.True(result1)

		result2 := u2.IsDeleted(s.DB)
		s.True(result2)

		// check if the user's sessions are deleted
		exist1 := s.RedisClient.Exists(context.Background(), sessionKey1).Val()
		s.Zero(exist1)
		exist2 := s.RedisClient.Exists(context.Background(), sessionKey2).Val()
		s.Zero(exist2)
	})
}

func (s *UserTestSuite) TestServiceIntegration_DisableUser() {
	s.Run("should disable user", func() {
		s.T().Parallel()

		sessionManager := session.NewRedisSessionManager(s.RedisClient)
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo, sessionManager)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)
		sessionKey := fmt.Sprintf("session:org:%v:user:%v", o.ID, u.ID)
		s.RedisClient.HSet(context.Background(), sessionKey, "jwt", gofakeit.UUID())

		err := svc.DisableUser(context.Background(), u.ID, "test")
		s.Require().NoError(err)

		result, err := svc.GetUserByID(context.Background(), u.ID)
		s.Require().NoError(err)
		s.NotNil(result.DisabledAt)
		s.Require().NotNil(result.Comment)
		s.Equal("test", *result.Comment)

		// check if the user's session is deleted
		exist := s.RedisClient.Exists(context.Background(), sessionKey).Val()
		s.Zero(exist)
	})
}

func (s *UserTestSuite) TestServiceIntegration_EnableUser() {
	s.Run("should enable user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		svc := user.NewService(repo, nil)
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
		svc := user.NewService(repo, nil)
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
		svc := user.NewService(repo, nil)
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
		svc := user.NewService(repo, nil)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserEmailNotVerified())

		err := svc.SetEmailVerified(context.Background(), u.ID)
		s.Require().NoError(err)

		result, err := svc.GetUserByID(context.Background(), u.ID)
		s.Require().NoError(err)
		s.True(result.IsEmailVerified)
	})
}
