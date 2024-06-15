package user_test

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/camelhr-api/internal/domains/session"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_GetUserByID(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		mockRepo.On("GetUserByID", context.Background(), int64(1)).
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByID(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return base.NotFoundError when user not found", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		mockRepo.On("GetUserByID", context.Background(), int64(1)).
			Return(user.User{}, sql.ErrNoRows)

		_, err := service.GetUserByID(context.Background(), int64(1))
		require.Error(t, err)
		require.IsType(t, &base.NotFoundError{}, err)
		assert.ErrorContains(t, err, "user not found for the given id")
	})

	t.Run("should return user by id", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		mockRepo.On("GetUserByAPIToken", context.Background(), "token").
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByAPIToken(context.Background(), "token")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when user not found for the given api token", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		mockRepo.On("GetUserByAPIToken", context.Background(), "token").
			Return(user.User{}, sql.ErrNoRows)

		_, err := service.GetUserByAPIToken(context.Background(), "token")
		require.Error(t, err)
		require.IsType(t, &base.NotFoundError{}, err)
		assert.ErrorContains(t, err, "user not found for the given api-token")
	})

	t.Run("should return user by api token", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

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

func TestService_GetUserByOrgSubdomainAPIToken(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		mockRepo.On("GetUserByOrgSubdomainAPIToken", context.Background(), "subdomain", "token").
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByOrgSubdomainAPIToken(context.Background(), "subdomain", "token")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when user not found for the given subdomain and api token", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		mockRepo.On("GetUserByOrgSubdomainAPIToken", context.Background(), "subdomain", "token").
			Return(user.User{}, sql.ErrNoRows)

		_, err := service.GetUserByOrgSubdomainAPIToken(context.Background(), "subdomain", "token")
		require.Error(t, err)
		require.IsType(t, &base.NotFoundError{}, err)
		assert.ErrorContains(t, err, "user not found for the given org-subdomain and api-token")
	})

	t.Run("should return error when subdomain is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		_, err := service.GetUserByOrgSubdomainAPIToken(context.Background(), "invalid_sub", "token")
		require.Error(t, err)
		assert.ErrorContains(t, err, "subdomain can only contain alphanumeric characters")
	})

	t.Run("should return user by organization subdomain and api token", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
		}

		mockRepo.On("GetUserByOrgSubdomainAPIToken", context.Background(), "subdomain", "token").
			Return(u, nil)

		result, err := service.GetUserByOrgSubdomainAPIToken(context.Background(), "subdomain", "token")
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestService_GetUserByOrgIDEmail(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		email := gofakeit.Email()

		mockRepo.On("GetUserByOrgIDEmail", context.Background(), int64(1), email).
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByOrgIDEmail(context.Background(), int64(1), email)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when user not found for the given organization id and email", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		email := gofakeit.Email()

		mockRepo.On("GetUserByOrgIDEmail", context.Background(), int64(1), email).
			Return(user.User{}, sql.ErrNoRows)

		_, err := service.GetUserByOrgIDEmail(context.Background(), int64(1), email)
		require.Error(t, err)
		require.IsType(t, &base.NotFoundError{}, err)
		assert.ErrorContains(t, err, "user not found for the given org-id and email")
	})

	t.Run("should return error when email is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		email := "invalid@invalid"

		_, err := service.GetUserByOrgIDEmail(context.Background(), int64(1), email)
		require.Error(t, err)
		assert.ErrorContains(t, err, "email must be a valid email address")
	})

	t.Run("should return user by organization id and email", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
		}

		mockRepo.On("GetUserByOrgIDEmail", context.Background(), u.ID, u.Email).
			Return(u, nil)

		result, err := service.GetUserByOrgIDEmail(context.Background(), u.ID, u.Email)
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestService_GetUserByOrgSubdomainEmail(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		email := gofakeit.Email()

		mockRepo.On("GetUserByOrgSubdomainEmail", context.Background(), "subdomain", email).
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByOrgSubdomainEmail(context.Background(), "subdomain", email)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when user not found for the given subdomain and email", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		email := gofakeit.Email()

		mockRepo.On("GetUserByOrgSubdomainEmail", context.Background(), "subdomain", email).
			Return(user.User{}, sql.ErrNoRows)

		_, err := service.GetUserByOrgSubdomainEmail(context.Background(), "subdomain", email)
		require.Error(t, err)
		require.IsType(t, &base.NotFoundError{}, err)
		assert.ErrorContains(t, err, "user not found for the given org-subdomain and email")
	})

	t.Run("should return error when subdomain is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		_, err := service.GetUserByOrgSubdomainEmail(context.Background(), "@#invalid", gofakeit.Email())
		require.Error(t, err)
		assert.ErrorContains(t, err, "subdomain can only contain alphanumeric characters")
	})

	t.Run("should return error when email is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		email := ""

		_, err := service.GetUserByOrgSubdomainEmail(context.Background(), "subdomain", email)
		require.Error(t, err)
		assert.ErrorContains(t, err, "email is required")
	})

	t.Run("should return user by organization subdomain and email", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		email := gofakeit.Email()

		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
		}

		mockRepo.On("GetUserByOrgSubdomainEmail", context.Background(), "subdomain", email).
			Return(u, nil)

		result, err := service.GetUserByOrgSubdomainEmail(context.Background(), "subdomain", email)
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestService_CreateUser(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		password := generatePassword()

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

	t.Run("should return error when email is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		password := generatePassword()

		_, err := service.CreateUser(context.Background(), int64(1), "invalid", password)
		require.Error(t, err)
		assert.ErrorContains(t, err, "email must be a valid email address")
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		_, err := service.CreateUser(context.Background(), int64(1), gofakeit.Email(), "invalid")
		require.Error(t, err)
		assert.ErrorContains(t, err, "password must be at least 8 characters in length")
	})

	t.Run("should create user", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		password := generatePassword()

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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		password := generatePassword()

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

	t.Run("should return error when email is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		password := generatePassword()

		_, err := service.CreateOwner(context.Background(), int64(1), "invalid", password)
		require.Error(t, err)
		assert.ErrorContains(t, err, "email must be a valid email address")
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		_, err := service.CreateOwner(context.Background(), int64(1), gofakeit.Email(), "invalid123")
		require.Error(t, err)
		assert.ErrorContains(t, err, "password must contain at least one uppercase letter")
	})

	t.Run("should create owner", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		password := generatePassword()

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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		password := generatePassword()

		mockRepo.On("ResetPassword", context.Background(), int64(1), fake.MockString).
			Return(assert.AnError)

		err := service.ResetPassword(context.Background(), int64(1), password)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		err := service.ResetPassword(context.Background(), int64(1), "Invalid123")
		require.Error(t, err)
		assert.ErrorContains(t, err, "password must contain at least one special character")
	})

	t.Run("should reset password", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		password := generatePassword()

		mockRepo.On("ResetPassword", context.Background(), int64(1), fake.MockString).
			Return(nil)

		err := service.ResetPassword(context.Background(), int64(1), password)
		require.NoError(t, err)
	})
}

func TestService_DeleteUser(t *testing.T) {
	t.Parallel()

	t.Run("should return error when comment is empty", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		err := service.DeleteUser(context.Background(), int64(1), "")
		require.Error(t, err)
		assert.True(t, base.IsInputValidationError(err))
		assert.ErrorContains(t, err, "comment is required")
	})

	t.Run("should return error when repo.GetUserByID return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		comment := gofakeit.Sentence(5)

		mockRepo.On("GetUserByID", context.Background(), int64(1)).
			Return(user.User{}, assert.AnError)

		err := service.DeleteUser(context.Background(), int64(1), comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when repo.DeleteUser return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		comment := gofakeit.Sentence(5)

		mockRepo.On("GetUserByID", context.Background(), int64(1)).
			Return(user.User{}, nil)
		mockRepo.On("DeleteUser", context.Background(), int64(1), comment).
			Return(assert.AnError)

		err := service.DeleteUser(context.Background(), int64(1), comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return base.NotFoundError when user not found for the given id", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		comment := gofakeit.Sentence(5)

		mockRepo.On("GetUserByID", context.Background(), int64(1)).
			Return(user.User{}, sql.ErrNoRows)

		err := service.DeleteUser(context.Background(), int64(1), comment)
		require.Error(t, err)
		require.IsType(t, &base.NotFoundError{}, err)
		assert.ErrorContains(t, err, "user not found for the given id")
	})

	t.Run("should return error when user is owner of the organization", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		comment := gofakeit.Sentence(5)

		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
			IsOwner:        true,
		}

		mockRepo.On("GetUserByID", context.Background(), u.ID).
			Return(u, nil)

		err := service.DeleteUser(context.Background(), u.ID, comment)
		require.Error(t, err)
		assert.ErrorIs(t, user.ErrUserIsOwner, err)
	})

	t.Run("should return error when user session deletion failed", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		sessionManager := session.NewMockSessionManager(t)
		service := user.NewService(mockRepo, sessionManager)
		comment := gofakeit.Sentence(5)

		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
		}

		mockRepo.On("GetUserByID", context.Background(), u.ID).
			Return(u, nil)
		mockRepo.On("DeleteUser", context.Background(), u.ID, comment).
			Return(nil)
		sessionManager.On("DeleteSession", context.Background(), u.ID, u.OrganizationID).
			Return(assert.AnError)

		err := service.DeleteUser(context.Background(), u.ID, comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should delete user", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		sessionManager := session.NewMockSessionManager(t)
		service := user.NewService(mockRepo, sessionManager)
		comment := gofakeit.Sentence(5)

		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
		}

		mockRepo.On("GetUserByID", context.Background(), u.ID).
			Return(u, nil)
		mockRepo.On("DeleteUser", context.Background(), u.ID, comment).
			Return(nil)
		sessionManager.On("DeleteSession", context.Background(), u.ID, u.OrganizationID).
			Return(nil)

		err := service.DeleteUser(context.Background(), u.ID, comment)
		require.NoError(t, err)
	})
}

func TestService_DisableUser(t *testing.T) {
	t.Parallel()

	t.Run("should return error when comment is empty", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		err := service.DisableUser(context.Background(), int64(1), "")
		require.Error(t, err)
		assert.True(t, base.IsInputValidationError(err))
		assert.ErrorContains(t, err, "comment is required")
	})

	t.Run("should return error when repo.GetUserByID return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		comment := gofakeit.SentenceSimple()

		mockRepo.On("GetUserByID", context.Background(), int64(1)).
			Return(user.User{}, assert.AnError)

		err := service.DisableUser(context.Background(), int64(1), comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when repo.DisableUser return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		comment := gofakeit.SentenceSimple()

		mockRepo.On("GetUserByID", context.Background(), int64(1)).
			Return(user.User{}, nil)
		mockRepo.On("DisableUser", context.Background(), int64(1), comment).
			Return(assert.AnError)

		err := service.DisableUser(context.Background(), int64(1), comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when user not found for the given id", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		comment := gofakeit.SentenceSimple()

		mockRepo.On("GetUserByID", context.Background(), int64(1)).
			Return(user.User{}, sql.ErrNoRows)

		err := service.DisableUser(context.Background(), int64(1), comment)
		require.Error(t, err)
		require.IsType(t, &base.NotFoundError{}, err)
		assert.ErrorContains(t, err, "user not found for the given id")
	})

	t.Run("should return error when user is owner of the organization", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		comment := gofakeit.SentenceSimple()

		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
			IsOwner:        true,
		}

		mockRepo.On("GetUserByID", context.Background(), u.ID).
			Return(u, nil)

		err := service.DisableUser(context.Background(), u.ID, comment)
		require.Error(t, err)
		assert.ErrorIs(t, user.ErrUserIsOwner, err)
	})

	t.Run("should return error when user session deletion failed", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		sessionManager := session.NewMockSessionManager(t)
		service := user.NewService(mockRepo, sessionManager)
		comment := gofakeit.SentenceSimple()

		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
		}

		mockRepo.On("GetUserByID", context.Background(), u.ID).
			Return(u, nil)
		mockRepo.On("DisableUser", context.Background(), u.ID, comment).
			Return(nil)
		sessionManager.On("DeleteSession", context.Background(), u.ID, u.OrganizationID).
			Return(assert.AnError)

		err := service.DisableUser(context.Background(), u.ID, comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should disable user", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		sessionManager := session.NewMockSessionManager(t)
		service := user.NewService(mockRepo, sessionManager)
		comment := gofakeit.SentenceSimple()

		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
			Email:          gofakeit.Email(),
		}

		mockRepo.On("GetUserByID", context.Background(), u.ID).
			Return(u, nil)
		mockRepo.On("DisableUser", context.Background(), u.ID, comment).
			Return(nil)
		sessionManager.On("DeleteSession", context.Background(), u.ID, u.OrganizationID).
			Return(nil)

		err := service.DisableUser(context.Background(), u.ID, comment)
		require.NoError(t, err)
	})
}

func TestService_EnableUser(t *testing.T) {
	t.Parallel()

	t.Run("should return error when comment is empty", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		err := service.DisableUser(context.Background(), int64(1), "")
		require.Error(t, err)
		assert.True(t, base.IsInputValidationError(err))
		assert.ErrorContains(t, err, "comment is required")
	})

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
		comment := gofakeit.SentenceSimple()

		mockRepo.On("EnableUser", context.Background(), int64(1), comment).
			Return(assert.AnError)

		err := service.EnableUser(context.Background(), int64(1), comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should enable user", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)
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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		mockRepo.On("GenerateAPIToken", context.Background(), int64(1)).
			Return(assert.AnError)

		err := service.GenerateAPIToken(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should generate api token", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		mockRepo.On("ResetAPIToken", context.Background(), int64(1)).
			Return(assert.AnError)

		err := service.ResetAPIToken(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should reset api token", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		mockRepo.On("SetEmailVerified", context.Background(), int64(1)).
			Return(assert.AnError)

		err := service.SetEmailVerified(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should set email verified", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo, nil)

		mockRepo.On("SetEmailVerified", context.Background(), int64(1)).
			Return(nil)

		err := service.SetEmailVerified(context.Background(), int64(1))
		require.NoError(t, err)
	})
}

// generatePassword generates a random password.
// It contains at least one lowercase letter, one uppercase letter, one special character, and one number.
// The minimum length of the password is 8 characters.
func generatePassword() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var (
		lowerChars    = "abcdefghijklmnopqrstuvwxyz"
		upperChars    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		specialChars  = "!@#$%^&*()-_=+,.?/:;{}[]`~"
		numberChars   = "0123456789"
		allChars      = lowerChars + upperChars + specialChars + numberChars
		password      = make([]byte, 0, 13)
		requiredChars = []byte{
			lowerChars[r.Intn(len(lowerChars))],
			upperChars[r.Intn(len(upperChars))],
			specialChars[r.Intn(len(specialChars))],
			numberChars[r.Intn(len(numberChars))],
		}
	)

	for range [9]int{} {
		password = append(
			password, allChars[r.Intn(len(allChars))])
	}

	password = append(password, requiredChars...)

	r.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
}
