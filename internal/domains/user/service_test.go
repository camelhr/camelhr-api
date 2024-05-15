package user_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_GetUserByID(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		mockRepo.On("GetUserByID", context.Background(), int64(1)).
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByID(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return user by id", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
		}

		mockRepo.On("GetUserByID", context.Background(), int64(1)).
			Return(u, nil)

		result, err := service.GetUserByID(context.Background(), int64(1))
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestService_GetUserByAPIToken(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		mockRepo.On("GetUserByAPIToken", context.Background(), "token").
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByAPIToken(context.Background(), "token")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return user by api token", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
		}

		mockRepo.On("GetUserByAPIToken", context.Background(), "token").
			Return(u, nil)

		result, err := service.GetUserByAPIToken(context.Background(), "token")
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestService_GetUserByOrgIDEmail(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		mockRepo.On("GetUserByOrgIDEmail", context.Background(), int64(1), "email").
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByOrgIDEmail(context.Background(), int64(1), "email")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return user by organization id and email", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
		}

		mockRepo.On("GetUserByOrgIDEmail", context.Background(), int64(1), "email").
			Return(u, nil)

		result, err := service.GetUserByOrgIDEmail(context.Background(), int64(1), "email")
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestService_GetUserByOrgSubdomainEmail(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		mockRepo.On("GetUserByOrgSubdomainEmail", context.Background(), "subdomain", "email").
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByOrgSubdomainEmail(context.Background(), "subdomain", "email")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return user by organization subdomain and email", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
		}

		mockRepo.On("GetUserByOrgSubdomainEmail", context.Background(), "subdomain", "email").
			Return(u, nil)

		result, err := service.GetUserByOrgSubdomainEmail(context.Background(), "subdomain", "email")
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestService_CreateUser(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)
		password := gofakeit.Password(true, true, true, true, false, 12)

		u := user.User{
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
			PasswordHash:   gofakeit.UUID(),
		}

		mockRepo.On("CreateUser", context.Background(), u.OrganizationID, u.Email, fake.MockString, false).
			Return(user.User{}, assert.AnError)

		_, err := service.CreateUser(context.Background(), u.OrganizationID, u.Email, password)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should create user", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)
		password := gofakeit.Password(true, true, true, true, false, 12)

		u := user.User{
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
			PasswordHash:   gofakeit.UUID(),
		}

		mockRepo.On("CreateUser", context.Background(), u.OrganizationID, u.Email, fake.MockString, false).
			Return(u, nil)

		result, err := service.CreateUser(context.Background(), u.OrganizationID, u.Email, password)
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestService_CreateOwner(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)
		password := gofakeit.Password(true, true, true, true, false, 12)

		u := user.User{
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
			PasswordHash:   gofakeit.UUID(),
		}

		mockRepo.On("CreateUser", context.Background(), u.OrganizationID, u.Email, fake.MockString, true).
			Return(user.User{}, assert.AnError)

		_, err := service.CreateOwner(context.Background(), u.OrganizationID, u.Email, password)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should create owner", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)
		password := gofakeit.Password(true, true, true, true, false, 12)

		u := user.User{
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
			PasswordHash:   gofakeit.UUID(),
		}

		mockRepo.On("CreateUser", context.Background(), u.OrganizationID, u.Email, fake.MockString, true).
			Return(u, nil)

		result, err := service.CreateOwner(context.Background(), u.OrganizationID, u.Email, password)
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestService_ResetPassword(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)
		password := gofakeit.Password(true, true, true, true, false, 12)

		mockRepo.On("ResetPassword", context.Background(), int64(1), fake.MockString).
			Return(assert.AnError)

		err := service.ResetPassword(context.Background(), int64(1), password)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should reset password", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)
		password := gofakeit.Password(true, true, true, true, false, 12)

		mockRepo.On("ResetPassword", context.Background(), int64(1), fake.MockString).
			Return(nil)

		err := service.ResetPassword(context.Background(), int64(1), password)
		require.NoError(t, err)
	})
}

func TestService_DeleteUser(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		mockRepo.On("DeleteUser", context.Background(), int64(1)).
			Return(assert.AnError)

		err := service.DeleteUser(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should delete user", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		mockRepo.On("DeleteUser", context.Background(), int64(1)).
			Return(nil)

		err := service.DeleteUser(context.Background(), int64(1))
		require.NoError(t, err)
	})
}

func TestService_DisableUser(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)
		comment := gofakeit.SentenceSimple()

		mockRepo.On("DisableUser", context.Background(), int64(1), comment).
			Return(assert.AnError)

		err := service.DisableUser(context.Background(), int64(1), comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should disable user", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)
		comment := gofakeit.SentenceSimple()

		mockRepo.On("DisableUser", context.Background(), int64(1), comment).
			Return(nil)

		err := service.DisableUser(context.Background(), int64(1), comment)
		require.NoError(t, err)
	})
}

func TestService_EnableUser(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)
		comment := gofakeit.SentenceSimple()

		mockRepo.On("EnableUser", context.Background(), int64(1), comment).
			Return(assert.AnError)

		err := service.EnableUser(context.Background(), int64(1), comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should enable user", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)
		comment := gofakeit.SentenceSimple()

		mockRepo.On("EnableUser", context.Background(), int64(1), comment).
			Return(nil)

		err := service.EnableUser(context.Background(), int64(1), comment)
		require.NoError(t, err)
	})
}

func TestService_GenerateAPIToken(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		mockRepo.On("GenerateAPIToken", context.Background(), int64(1)).
			Return(assert.AnError)

		err := service.GenerateAPIToken(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should generate api token", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		mockRepo.On("GenerateAPIToken", context.Background(), int64(1)).
			Return(nil)

		err := service.GenerateAPIToken(context.Background(), int64(1))
		require.NoError(t, err)
	})
}

func TestService_ResetAPIToken(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		mockRepo.On("ResetAPIToken", context.Background(), int64(1)).
			Return(assert.AnError)

		err := service.ResetAPIToken(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should reset api token", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		mockRepo.On("ResetAPIToken", context.Background(), int64(1)).
			Return(nil)

		err := service.ResetAPIToken(context.Background(), int64(1))
		require.NoError(t, err)
	})
}

func TestService_SetEmailVerified(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		mockRepo.On("SetEmailVerified", context.Background(), int64(1)).
			Return(assert.AnError)

		err := service.SetEmailVerified(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should set email verified", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewRepositoryMock(t)
		service := user.NewService(mockRepo)

		mockRepo.On("SetEmailVerified", context.Background(), int64(1)).
			Return(nil)

		err := service.SetEmailVerified(context.Background(), int64(1))
		require.NoError(t, err)
	})
}
