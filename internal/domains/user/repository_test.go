package user_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetUserByID(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Get", context.Background(), mock.Anything, tests.QueryMatcher("getUserByIDQuery"), int64(1)).
			Return(assert.AnError)

		_, err := repo.GetUserByID(context.Background(), 1)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return a user", func(t *testing.T) {
		t.Parallel()

		var emptyUser user.User

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		u := user.User{ID: 1}

		mockDB.On("Get", context.Background(), &emptyUser, tests.QueryMatcher("getUserByIDQuery"), int64(1)).
			Run(func(args mock.Arguments) {
				// populate the passed argument with the user
				arg, ok := args.Get(1).(*user.User)
				require.True(t, ok)
				*arg = u
			}).Return(nil)

		result, err := repo.GetUserByID(context.Background(), 1)
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestRepository_GetUserByAPIToken(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Get", context.Background(), mock.Anything, tests.QueryMatcher("getUserByAPITokenQuery"), "token").
			Return(assert.AnError)

		_, err := repo.GetUserByAPIToken(context.Background(), "token")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return a user", func(t *testing.T) {
		t.Parallel()

		var emptyUser user.User

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		u := user.User{ID: 1}

		mockDB.On("Get", context.Background(), &emptyUser, tests.QueryMatcher("getUserByAPITokenQuery"), "token").
			Run(func(args mock.Arguments) {
				// populate the passed argument with the user
				arg, ok := args.Get(1).(*user.User)
				require.True(t, ok)
				*arg = u
			}).Return(nil)

		result, err := repo.GetUserByAPIToken(context.Background(), "token")
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestRepository_GetUserByOrgSubdomainAPIToken(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Get", context.Background(), mock.Anything, tests.QueryMatcher("getUserByOrgSubdomainAPITokenQuery"),
			"subdomain", "token").Return(assert.AnError)

		_, err := repo.GetUserByOrgSubdomainAPIToken(context.Background(), "subdomain", "token")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return a user", func(t *testing.T) {
		t.Parallel()

		var emptyUser user.User

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		u := user.User{ID: 1}

		mockDB.On("Get", context.Background(), &emptyUser, tests.QueryMatcher("getUserByOrgSubdomainAPITokenQuery"),
			"subdomain", "token").Run(func(args mock.Arguments) {
			// populate the passed argument with the user
			arg, ok := args.Get(1).(*user.User)
			require.True(t, ok)
			*arg = u
		}).Return(nil)

		result, err := repo.GetUserByOrgSubdomainAPIToken(context.Background(), "subdomain", "token")
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestRepository_GetUserByOrgIDEmail(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Get", context.Background(), mock.Anything,
			tests.QueryMatcher("getUserByOrgIDEmailQuery"), int64(1), "email").
			Return(assert.AnError)

		_, err := repo.GetUserByOrgIDEmail(context.Background(), 1, "email")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return a user", func(t *testing.T) {
		t.Parallel()

		var emptyUser user.User

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		u := user.User{ID: 1}

		mockDB.On("Get", context.Background(), &emptyUser, tests.QueryMatcher("getUserByOrgIDEmailQuery"), int64(1), "email").
			Run(func(args mock.Arguments) {
				// populate the passed argument with the user
				arg, ok := args.Get(1).(*user.User)
				require.True(t, ok)
				*arg = u
			}).Return(nil)

		result, err := repo.GetUserByOrgIDEmail(context.Background(), 1, "email")
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestRepository_GetUserByOrgSubdomainEmail(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Get", context.Background(), mock.Anything, tests.QueryMatcher("getUserByOrgSubdomainEmailQuery"),
			"subdomain", "email").Return(assert.AnError)

		_, err := repo.GetUserByOrgSubdomainEmail(context.Background(), "subdomain", "email")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return a user", func(t *testing.T) {
		t.Parallel()

		var emptyUser user.User

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		u := user.User{ID: 1}

		mockDB.On("Get", context.Background(), &emptyUser,
			tests.QueryMatcher("getUserByOrgSubdomainEmailQuery"), "subdomain", "email").
			Run(func(args mock.Arguments) {
				// populate the passed argument with the user
				arg, ok := args.Get(1).(*user.User)
				require.True(t, ok)
				*arg = u
			}).Return(nil)

		result, err := repo.GetUserByOrgSubdomainEmail(context.Background(), "subdomain", "email")
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestRepository_CreateUser(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), mock.Anything,
			tests.QueryMatcher("createUserQuery"), int64(1), "email", "password", true).Return(assert.AnError)

		_, err := repo.CreateUser(context.Background(), 1, "email", "password", true)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return a user", func(t *testing.T) {
		t.Parallel()

		var emptyUser user.User

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		u := user.User{ID: 1}

		mockDB.On("Exec", context.Background(), &emptyUser,
			tests.QueryMatcher("createUserQuery"), int64(1), "email", "password", false).
			Run(func(args mock.Arguments) {
				// populate the passed argument with the user
				arg, ok := args.Get(1).(*user.User)
				require.True(t, ok)
				*arg = u
			}).Return(nil)

		result, err := repo.CreateUser(context.Background(), 1, "email", "password", false)
		require.NoError(t, err)
		assert.Equal(t, u, result)
	})
}

func TestRepository_ResetPassword(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("resetPasswordQuery"), int64(1), "password").
			Return(assert.AnError)

		err := repo.ResetPassword(context.Background(), 1, "password")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the password is reset", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("resetPasswordQuery"), int64(1), "password").
			Return(nil)

		err := repo.ResetPassword(context.Background(), 1, "password")
		require.NoError(t, err)
	})
}

func TestRepository_DeleteUser(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("deleteUserQuery"), int64(1)).
			Return(assert.AnError)

		err := repo.DeleteUser(context.Background(), 1)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when user is deleted", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("deleteUserQuery"), int64(1)).
			Return(nil)

		err := repo.DeleteUser(context.Background(), 1)
		require.NoError(t, err)
	})
}

func TestRepository_DisableUser(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)
		comment := gofakeit.SentenceSimple()

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("disableUserQuery"), int64(1), comment).
			Return(assert.AnError)

		err := repo.DisableUser(context.Background(), 1, comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when user is disabled", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)
		comment := gofakeit.SentenceSimple()

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("disableUserQuery"), int64(1), comment).
			Return(nil)

		err := repo.DisableUser(context.Background(), 1, comment)
		require.NoError(t, err)
	})
}

func TestRepository_EnableUser(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)
		comment := gofakeit.SentenceSimple()

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("enableUserQuery"), int64(1), comment).
			Return(assert.AnError)

		err := repo.EnableUser(context.Background(), 1, comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when user is enabled", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)
		comment := gofakeit.SentenceSimple()

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("enableUserQuery"), int64(1), comment).
			Return(nil)

		err := repo.EnableUser(context.Background(), 1, comment)
		require.NoError(t, err)
	})
}

func TestRepository_GenerateAPIToken(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("generateAPITokenQuery"), int64(1)).
			Return(assert.AnError)

		err := repo.GenerateAPIToken(context.Background(), 1)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when api token is generated", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("generateAPITokenQuery"), int64(1)).
			Return(nil)

		err := repo.GenerateAPIToken(context.Background(), 1)
		require.NoError(t, err)
	})
}

func TestRepository_ResetAPIToken(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("resetAPITokenQuery"), int64(1)).
			Return(assert.AnError)

		err := repo.ResetAPIToken(context.Background(), 1)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when api token is reset", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("resetAPITokenQuery"), int64(1)).
			Return(nil)

		err := repo.ResetAPIToken(context.Background(), 1)
		require.NoError(t, err)
	})
}

func TestRepository_SetEmailVerified(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("setEmailVerifiedQuery"), int64(1)).
			Return(assert.AnError)

		err := repo.SetEmailVerified(context.Background(), 1)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when email is verified", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := user.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("setEmailVerifiedQuery"), int64(1)).
			Return(nil)

		err := repo.SetEmailVerified(context.Background(), 1)
		require.NoError(t, err)
	})
}
