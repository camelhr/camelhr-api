package user_test

import (
	"context"
	"database/sql"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
)

func (s *UserTestSuite) TestRepositoryIntegration_Triggers() {
	// tests to ensure that the triggers are working as expected
	s.Run("should forbid truncate operation on users table", func() {
		s.T().Parallel()
		err := s.DB.Exec(context.Background(), nil, "TRUNCATE users CASCADE")
		s.Require().Error(err)
		s.ErrorContains(err, "TRUNCATE operation on table users is not allowed: prevent_truncate_on_users")
	})

	s.Run("should forbid delete operation on users table", func() {
		s.T().Parallel()
		fake.NewUser(s.DB, fake.NewOrganization(s.DB).ID)
		err := s.DB.Exec(context.Background(), nil, "DELETE FROM users WHERE user_id = 1")
		s.Require().Error(err)
		s.ErrorContains(err, "DELETE operation on table users is not allowed: prevent_hard_delete_on_users")
	})

	s.Run("should soft delete users if the organization is deleted", func() {
		s.T().Parallel()
		o := fake.NewOrganization(s.DB)
		u1 := fake.NewUser(s.DB, o.ID)
		u2 := fake.NewUser(s.DB, o.ID)
		u3 := fake.NewUser(s.DB, o.ID, fake.UserDeleted())
		o.Delete(s.DB)

		u1Latest := u1.FetchLatest(s.DB)
		s.NotNil(u1Latest.DeletedAt)
		s.NotNil(u1Latest.Comment)
		s.Equal("deletion_reason: associated_organization_deleted", *u1Latest.Comment)

		u2Latest := u2.FetchLatest(s.DB)
		s.NotNil(u2Latest.DeletedAt)
		s.NotNil(u2Latest.Comment)
		s.Equal("deletion_reason: associated_organization_deleted", *u2Latest.Comment)

		// should not delete already deleted user
		u3Latest := u3.FetchLatest(s.DB)
		s.NotNil(u3Latest.DeletedAt)
		s.Nil(u3Latest.Comment)
		s.Equal(u3.DeletedAt, u3Latest.DeletedAt)
	})
}

func (s *UserTestSuite) TestRepositoryIntegration_GetUserByID() {
	s.Run("should return user by id", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		result, err := repo.GetUserByID(context.Background(), u.ID)
		s.Require().NoError(err)
		s.Equal(u.User, result)
	})

	s.Run("should return error when user does not exist", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)

		_, err := repo.GetUserByID(context.Background(), 0)
		s.Require().Error(err)
		s.ErrorIs(err, sql.ErrNoRows)
	})

	s.Run("should return error when user is deleted", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		u := fake.NewUser(s.DB, fake.NewOrganization(s.DB).ID, fake.UserDeleted())

		_, err := repo.GetUserByID(context.Background(), u.ID)
		s.Require().Error(err)
		s.ErrorIs(err, sql.ErrNoRows)
	})
}

func (s *UserTestSuite) TestRepositoryIntegration_GetUserByAPIToken() {
	s.Run("should return user by api token", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)
		s.Require().NotNil(u.APIToken)

		result, err := repo.GetUserByAPIToken(context.Background(), *u.APIToken)
		s.Require().NoError(err)
		s.Equal(u.User, result)
	})

	s.Run("should return error when user does not exist for api token", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)

		_, err := repo.GetUserByAPIToken(context.Background(), gofakeit.UUID())
		s.Require().Error(err)
		s.ErrorIs(err, sql.ErrNoRows)
	})

	s.Run("should return error when user is deleted", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		u := fake.NewUser(s.DB, fake.NewOrganization(s.DB).ID, fake.UserDeleted())
		s.Require().NotNil(u.APIToken)

		_, err := repo.GetUserByAPIToken(context.Background(), *u.APIToken)
		s.Require().Error(err)
		s.ErrorIs(err, sql.ErrNoRows)
	})
}

func (s *UserTestSuite) TestRepositoryIntegration_GetUserByOrgSubdomainAPIToken() {
	s.Run("should return user by org subdomain and api token", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)
		s.Require().NotNil(u.APIToken)

		result, err := repo.GetUserByOrgSubdomainAPIToken(context.Background(), o.Subdomain, *u.APIToken)
		s.Require().NoError(err)
		s.Equal(u.User, result)
	})

	s.Run("should return error when user does not exist for api token", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)

		_, err := repo.GetUserByOrgSubdomainAPIToken(context.Background(), o.Subdomain, gofakeit.UUID())
		s.Require().Error(err)
		s.ErrorIs(err, sql.ErrNoRows)
	})

	s.Run("should return error when user does not exist for subdomain", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o1 := fake.NewOrganization(s.DB)
		o2 := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o1.ID)
		s.Require().NotNil(u.APIToken)

		_, err := repo.GetUserByOrgSubdomainAPIToken(context.Background(), o2.Subdomain, *u.APIToken)
		s.Require().Error(err)
		s.ErrorIs(err, sql.ErrNoRows)
	})

	s.Run("should return error when user is deleted", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDeleted())
		s.Require().NotNil(u.APIToken)

		_, err := repo.GetUserByOrgSubdomainAPIToken(context.Background(), o.Subdomain, *u.APIToken)
		s.Require().Error(err)
		s.ErrorIs(err, sql.ErrNoRows)
	})
}

func (s *UserTestSuite) TestRepositoryIntegration_GetUserByOrgIDEmail() {
	s.Run("should return user by org id and email", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		result, err := repo.GetUserByOrgIDEmail(context.Background(), o.ID, u.Email)
		s.Require().NoError(err)
		s.Equal(u.User, result)
	})

	s.Run("should return error when user does not exist", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)

		_, err := repo.GetUserByOrgIDEmail(context.Background(), 0, gofakeit.Email())
		s.Require().Error(err)
		s.ErrorIs(err, sql.ErrNoRows)
	})

	s.Run("should return error when user is deleted", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		u := fake.NewUser(s.DB, fake.NewOrganization(s.DB).ID, fake.UserDeleted())

		_, err := repo.GetUserByOrgIDEmail(context.Background(), u.OrganizationID, u.Email)
		s.Require().Error(err)
		s.ErrorIs(err, sql.ErrNoRows)
	})
}

func (s *UserTestSuite) TestRepositoryIntegration_GetUserByOrgSubdomainEmail() {
	s.Run("should return user by org subdomain and email", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		result, err := repo.GetUserByOrgSubdomainEmail(context.Background(), o.Subdomain, u.Email)
		s.Require().NoError(err)
		s.Equal(u.User, result)
	})

	s.Run("should return error when user does not exist", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)

		_, err := repo.GetUserByOrgSubdomainEmail(context.Background(), "invalid-subdomain", gofakeit.Email())
		s.Require().Error(err)
		s.ErrorIs(err, sql.ErrNoRows)
	})

	s.Run("should return error when user is deleted", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		u := fake.NewUser(s.DB, fake.NewOrganization(s.DB).ID, fake.UserDeleted())

		_, err := repo.GetUserByOrgSubdomainEmail(context.Background(), "invalid-subdomain", u.Email)
		s.Require().Error(err)
		s.ErrorIs(err, sql.ErrNoRows)
	})
}

func (s *UserTestSuite) TestRepositoryIntegration_CreateUser() {
	s.Run("should create user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := user.User{
			OrganizationID: o.ID,
			Email:          gofakeit.Email(),
			PasswordHash:   gofakeit.UUID(),
		}

		result, err := repo.CreateUser(context.Background(), u.OrganizationID, u.Email, u.PasswordHash, false)
		s.Require().NoError(err)
		s.Equal(u.OrganizationID, result.OrganizationID)
		s.Equal(u.Email, result.Email)
		s.Equal(u.PasswordHash, result.PasswordHash)
		s.False(result.IsOwner)
		s.False(result.IsEmailVerified) // email is not verified by default
		s.Nil(result.APIToken)
		s.Nil(result.DisabledAt)
		s.Nil(result.Comment)
		s.NotZero(result.CreatedAt)
		s.NotZero(result.UpdatedAt)
		s.Equal(result.CreatedAt, result.UpdatedAt)
		s.Equal(time.UTC, result.CreatedAt.Location())
		s.Equal(time.UTC, result.UpdatedAt.Location())
		s.WithinDuration(time.Now().UTC(), result.CreatedAt, 1*time.Minute)
		s.WithinDuration(time.Now().UTC(), result.UpdatedAt, 1*time.Minute)
		s.Nil(result.DeletedAt)
	})

	s.Run("should not create user with existing email", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		_, err := repo.CreateUser(context.Background(), o.ID, u.Email, gofakeit.UUID(), false)
		s.Require().Error(err)
		s.ErrorContains(err, "duplicate key value violates unique constraint")
	})

	s.Run("should not create user with invalid organization id", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		u := user.User{
			OrganizationID: 0,
			Email:          gofakeit.Email(),
			PasswordHash:   gofakeit.UUID(),
		}

		_, err := repo.CreateUser(context.Background(), u.OrganizationID, u.Email, u.PasswordHash, false)
		s.Require().Error(err)
		s.ErrorContains(err, "insert or update on table \"users\" violates foreign key constraint")
	})

	s.Run("should not create user with invalid email", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := user.User{
			OrganizationID: o.ID,
			Email:          "invalid-email",
			PasswordHash:   gofakeit.UUID(),
		}

		_, err := repo.CreateUser(context.Background(), u.OrganizationID, u.Email, u.PasswordHash, false)
		s.Require().Error(err)
		s.ErrorContains(err, "new row for relation \"users\" violates check constraint \"users_email_check\"")
	})

	s.Run("should not create more than one owner", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		fake.NewUser(s.DB, o.ID, fake.UserIsOwner())
		owner2 := user.User{
			OrganizationID: o.ID,
			Email:          gofakeit.Email(),
			PasswordHash:   gofakeit.UUID(),
		}

		_, err := repo.CreateUser(context.Background(), owner2.OrganizationID, owner2.Email, owner2.PasswordHash, true)
		s.Require().Error(err)
		s.ErrorContains(err, "duplicate key value violates unique constraint \"idx_users_owner_per_organization\"")
	})
}

func (s *UserTestSuite) TestRepositoryIntegration_ResetPassword() {
	s.Run("should reset password", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)
		newPasswordHash := gofakeit.UUID()

		err := repo.ResetPassword(context.Background(), u.ID, newPasswordHash)
		s.Require().NoError(err)

		result, err := repo.GetUserByID(context.Background(), u.ID)
		s.Require().NoError(err)
		s.Equal(newPasswordHash, result.PasswordHash)
		s.Equal(time.UTC, result.UpdatedAt.Location())
		s.WithinDuration(time.Now().UTC(), result.UpdatedAt, 1*time.Minute)
		s.GreaterOrEqual(result.UpdatedAt.Unix(), result.CreatedAt.Unix())
	})
}

func (s *UserTestSuite) TestRepositoryIntegration_DeleteUser() {
	s.Run("should delete user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		err := repo.DeleteUser(context.Background(), u.ID)
		s.Require().NoError(err)

		isDeleted := u.IsDeleted(s.DB)
		s.True(isDeleted)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.NotNil(result.DeletedAt)
		s.Equal(time.UTC, result.DeletedAt.Location())
		s.WithinDuration(time.Now().UTC(), *result.DeletedAt, 1*time.Minute)
	})

	s.Run("should not delete user if already deleted", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDeleted())

		err := repo.DeleteUser(context.Background(), u.ID)
		s.Require().NoError(err)

		isDeleted := u.IsDeleted(s.DB)
		s.True(isDeleted)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.Equal(u.DeletedAt, result.DeletedAt)
	})

	s.Run("should not delete if user is owner", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserIsOwner())

		err := repo.DeleteUser(context.Background(), u.ID)
		s.Require().NoError(err)

		isDeleted := u.IsDeleted(s.DB)
		s.False(isDeleted)
	})
}

func (s *UserTestSuite) TestRepositoryIntegration_DisableUser() {
	s.Run("should disable user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)
		comment := gofakeit.SentenceSimple()

		err := repo.DisableUser(context.Background(), u.ID, comment)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.NotNil(result.DisabledAt)
		s.Equal(time.UTC, result.DisabledAt.Location())
		s.WithinDuration(time.Now().UTC(), *result.DisabledAt, 1*time.Minute)
		s.GreaterOrEqual(result.DisabledAt.Unix(), result.CreatedAt.Unix())
		s.NotNil(result.Comment)
		s.Equal(comment, *result.Comment)
	})

	s.Run("should not disable user if already disabled", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDisabled())
		comment := gofakeit.SentenceSimple()

		err := repo.DisableUser(context.Background(), u.ID, comment)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.Require().NotNil(result.DisabledAt)
		s.Equal(u.DisabledAt, result.DisabledAt)
	})

	s.Run("should not disable deleted user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDeleted())
		comment := gofakeit.SentenceSimple()

		err := repo.DisableUser(context.Background(), u.ID, comment)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.Require().NotNil(result.DeletedAt)
		s.Nil(result.DisabledAt)
	})

	s.Run("should not disable user if user is owner", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserIsOwner())
		comment := gofakeit.SentenceSimple()

		err := repo.DisableUser(context.Background(), u.ID, comment)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.Nil(result.DisabledAt)
	})
}

func (s *UserTestSuite) TestRepositoryIntegration_EnableUser() {
	s.Run("should enable user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDisabled())
		comment := gofakeit.SentenceSimple()

		err := repo.EnableUser(context.Background(), u.ID, comment)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.Nil(result.DisabledAt)
		s.NotNil(result.Comment)
		s.Equal(comment, *result.Comment)
	})
}

func (s *UserTestSuite) TestRepositoryIntegration_GenerateAPIToken() {
	s.Run("should generate api token", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)

		err := repo.GenerateAPIToken(context.Background(), u.ID)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.NotNil(result.APIToken)
		s.Equal(time.UTC, result.UpdatedAt.Location())
		s.WithinDuration(time.Now().UTC(), result.UpdatedAt, 1*time.Minute)
		s.GreaterOrEqual(result.UpdatedAt.Unix(), result.CreatedAt.Unix())
	})

	s.Run("should not generate api token if already exists", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)
		s.Require().NotNil(u.APIToken)

		err := repo.GenerateAPIToken(context.Background(), u.ID)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.NotNil(result.APIToken)
		s.Equal(u.APIToken, result.APIToken)
	})

	s.Run("should not generate api token for deleted user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDeleted(), fake.UserWithoutToken())

		err := repo.GenerateAPIToken(context.Background(), u.ID)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.Nil(result.APIToken)
	})
}

func (s *UserTestSuite) TestRepositoryIntegration_ResetAPIToken() {
	s.Run("should reset api token", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID)
		s.Require().NotNil(u.APIToken)

		err := repo.ResetAPIToken(context.Background(), u.ID)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.NotNil(result.APIToken)
		s.NotEqual(u.APIToken, result.APIToken)
		s.Equal(time.UTC, result.UpdatedAt.Location())
		s.WithinDuration(time.Now().UTC(), result.UpdatedAt, 1*time.Minute)
		s.GreaterOrEqual(result.UpdatedAt.Unix(), result.CreatedAt.Unix())
	})

	s.Run("should not reset api token if not exists", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserWithoutToken())

		err := repo.ResetAPIToken(context.Background(), u.ID)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.Nil(result.APIToken)
	})

	s.Run("should not reset api token for deleted user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDeleted())

		err := repo.ResetAPIToken(context.Background(), u.ID)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.Require().NotNil(result.DeletedAt)
		s.Equal(u.APIToken, result.APIToken)
	})

	s.Run("should not reset api token for disabled user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDisabled())

		err := repo.ResetAPIToken(context.Background(), u.ID)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.Require().NotNil(result.DisabledAt)
		s.Equal(u.APIToken, result.APIToken)
	})
}

func (s *UserTestSuite) TestRepositoryIntegration_SetEmailVerified() {
	s.Run("should set email verified flag", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserEmailNotVerified())

		err := repo.SetEmailVerified(context.Background(), u.ID)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.True(result.IsEmailVerified)
		s.Equal(time.UTC, result.UpdatedAt.Location())
		s.WithinDuration(time.Now().UTC(), result.UpdatedAt, 1*time.Minute)
		s.GreaterOrEqual(result.UpdatedAt.Unix(), result.CreatedAt.Unix())
	})

	s.Run("should not set email verified flag for deleted user", func() {
		s.T().Parallel()
		repo := user.NewRepository(s.DB)
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserDeleted(), fake.UserEmailNotVerified())

		err := repo.SetEmailVerified(context.Background(), u.ID)
		s.Require().NoError(err)

		result := u.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.Require().NotNil(result.DeletedAt)
		s.False(result.IsEmailVerified)
	})
}
