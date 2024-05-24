package user_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

		mockRepo.On("GetUserByID", context.Background(), int64(1)).
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByID(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return user by id", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

		mockRepo.On("GetUserByAPIToken", context.Background(), "token").
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByAPIToken(context.Background(), "token")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return user by api token", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
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

func TestService_GetUserByOrgSubdomainAPIToken(t *testing.T) {
	t.Parallel()

	t.Run("should return error when repository return error", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

		mockRepo.On("GetUserByOrgSubdomainAPIToken", context.Background(), "subdomain", "token").
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByOrgSubdomainAPIToken(context.Background(), "subdomain", "token")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when subdomain is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

		_, err := service.GetUserByOrgSubdomainAPIToken(context.Background(), "invalid_sub", "token")
		require.Error(t, err)
		assert.ErrorContains(t, err, "subdomain can only contain alphanumeric characters")
	})

	t.Run("should return user by organization subdomain and api token", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

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
		service := user.NewService(mockRepo)
		email := gofakeit.Email()

		mockRepo.On("GetUserByOrgIDEmail", context.Background(), int64(1), email).
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByOrgIDEmail(context.Background(), int64(1), email)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when email is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)
		email := "invalid@invalid"

		_, err := service.GetUserByOrgIDEmail(context.Background(), int64(1), email)
		require.Error(t, err)
		assert.ErrorContains(t, err, "email must be a valid email address")
	})

	t.Run("should return user by organization id and email", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

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
		service := user.NewService(mockRepo)
		email := gofakeit.Email()

		mockRepo.On("GetUserByOrgSubdomainEmail", context.Background(), "subdomain", email).
			Return(user.User{}, assert.AnError)

		_, err := service.GetUserByOrgSubdomainEmail(context.Background(), "subdomain", email)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when subdomain is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

		_, err := service.GetUserByOrgSubdomainEmail(context.Background(), "@#invalid", gofakeit.Email())
		require.Error(t, err)
		assert.ErrorContains(t, err, "subdomain can only contain alphanumeric characters")
	})

	t.Run("should return error when email is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)
		email := ""

		_, err := service.GetUserByOrgSubdomainEmail(context.Background(), "subdomain", email)
		require.Error(t, err)
		assert.ErrorContains(t, err, "email is required")
	})

	t.Run("should return user by organization subdomain and email", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)
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
		service := user.NewService(mockRepo)
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
		service := user.NewService(mockRepo)
		password := generatePassword()

		_, err := service.CreateUser(context.Background(), int64(1), "invalid", password)
		require.Error(t, err)
		assert.ErrorContains(t, err, "email must be a valid email address")
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

		_, err := service.CreateUser(context.Background(), int64(1), gofakeit.Email(), "invalid")
		require.Error(t, err)
		assert.ErrorContains(t, err, "password must be at least 8 characters in length")
	})

	t.Run("should create user", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)
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
		service := user.NewService(mockRepo)
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
		service := user.NewService(mockRepo)
		password := generatePassword()

		_, err := service.CreateOwner(context.Background(), int64(1), "invalid", password)
		require.Error(t, err)
		assert.ErrorContains(t, err, "email must be a valid email address")
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

		_, err := service.CreateOwner(context.Background(), int64(1), gofakeit.Email(), "invalid123")
		require.Error(t, err)
		assert.ErrorContains(t, err, "password must contain at least one uppercase letter")
	})

	t.Run("should create owner", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)
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
		service := user.NewService(mockRepo)
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
		service := user.NewService(mockRepo)

		err := service.ResetPassword(context.Background(), int64(1), "Invalid123")
		require.Error(t, err)
		assert.ErrorContains(t, err, "password must contain at least one special character")
	})

	t.Run("should reset password", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)
		password := generatePassword()

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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

		mockRepo.On("DeleteUser", context.Background(), int64(1)).
			Return(assert.AnError)

		err := service.DeleteUser(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should delete user", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)
		comment := gofakeit.SentenceSimple()

		mockRepo.On("DisableUser", context.Background(), int64(1), comment).
			Return(assert.AnError)

		err := service.DisableUser(context.Background(), int64(1), comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when comment is empty", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

		err := service.DisableUser(context.Background(), int64(1), "")
		require.Error(t, err)
		assert.ErrorContains(t, err, "comment is required")
	})

	t.Run("should disable user", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)
		comment := gofakeit.SentenceSimple()

		mockRepo.On("EnableUser", context.Background(), int64(1), comment).
			Return(assert.AnError)

		err := service.EnableUser(context.Background(), int64(1), comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when comment is empty", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

		err := service.DisableUser(context.Background(), int64(1), "")
		require.Error(t, err)
		assert.ErrorContains(t, err, "comment is required")
	})

	t.Run("should enable user", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

		mockRepo.On("GenerateAPIToken", context.Background(), int64(1)).
			Return(assert.AnError)

		err := service.GenerateAPIToken(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should generate api token", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

		mockRepo.On("ResetAPIToken", context.Background(), int64(1)).
			Return(assert.AnError)

		err := service.ResetAPIToken(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should reset api token", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
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

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

		mockRepo.On("SetEmailVerified", context.Background(), int64(1)).
			Return(assert.AnError)

		err := service.SetEmailVerified(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should set email verified", func(t *testing.T) {
		t.Parallel()

		mockRepo := user.NewMockRepository(t)
		service := user.NewService(mockRepo)

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
		password      []byte
		requiredChars = []byte{
			lowerChars[r.Intn(len(lowerChars))],
			upperChars[r.Intn(len(upperChars))],
			specialChars[r.Intn(len(specialChars))],
			numberChars[r.Intn(len(numberChars))],
		}
	)

	for i := 0; i < 9; i++ {
		password = append(password, allChars[r.Intn(len(allChars))])
	}

	password = append(password, requiredChars...)

	r.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
}
