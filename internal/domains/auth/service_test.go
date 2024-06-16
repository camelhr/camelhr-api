package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/camelhr/camelhr-api/internal/domains/auth"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/domains/session"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const (
	validPassword = "@paSSw0rd"
)

func TestService_Register(t *testing.T) {
	t.Parallel()

	t.Run("should return error when subdomain already exists", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		email := gofakeit.Email()
		orgName := gofakeit.Company()
		subdomain := gofakeit.LetterN(30)

		orgService := organization.NewMockService(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).
			Return(organization.Organization{ID: gofakeit.Int64(), Subdomain: subdomain}, nil)

		authService := auth.NewService("", nil, orgService, nil, nil)
		err := authService.Register(ctx, email, validPassword, subdomain, orgName)

		require.Error(t, err)
		require.ErrorIs(t, auth.ErrSubdomainAlreadyExists, err)
	})

	t.Run("should return error when orgService.GetOrganizationBySubdomain returns error other than not found",
		func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			email := gofakeit.Email()

			orgName := gofakeit.Company()
			subdomain := gofakeit.LetterN(30)

			orgService := organization.NewMockService(t)
			orgService.On("GetOrganizationBySubdomain", ctx, subdomain).
				Return(organization.Organization{}, assert.AnError)

			authService := auth.NewService("", nil, orgService, nil, nil)
			err := authService.Register(ctx, email, validPassword, subdomain, orgName)

			require.Error(t, err)
			require.ErrorIs(t, assert.AnError, err)
		})

	t.Run("should return error when transactor.WithTx returns error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		email := gofakeit.Email()
		orgName := gofakeit.Company()
		subdomain := gofakeit.LetterN(30)
		notFoundErr := base.NewNotFoundError("not found")

		orgService := organization.NewMockService(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).
			Return(organization.Organization{}, notFoundErr)

		transactor := database.NewMockTransactor(t)
		transactor.On("WithTx", ctx, mock.Anything).Return(assert.AnError)

		authService := auth.NewService("", transactor, orgService, nil, nil)
		err := authService.Register(ctx, email, validPassword, subdomain, orgName)

		require.Error(t, err)
		require.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should not return error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		email := gofakeit.Email()
		orgName := gofakeit.Company()
		subdomain := gofakeit.LetterN(30)
		notFoundErr := base.NewNotFoundError("not found")

		orgService := organization.NewMockService(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).Return(organization.Organization{}, notFoundErr)

		transactor := database.NewMockTransactor(t)
		transactor.On("WithTx", ctx, mock.Anything).Return(nil)

		authService := auth.NewService("", transactor, orgService, nil, nil)
		err := authService.Register(ctx, email, validPassword, subdomain, orgName)

		require.NoError(t, err)
	})
}

func TestService_Login(t *testing.T) {
	t.Parallel()

	t.Run("should return error when orgService.GetOrganizationBySubdomain returns error",
		func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			subdomain := gofakeit.LetterN(30)

			orgService := organization.NewMockService(t)
			orgService.On("GetOrganizationBySubdomain", ctx, subdomain).Return(organization.Organization{}, assert.AnError)

			authService := auth.NewService("secret", nil, orgService, nil, nil)
			_, _, err := authService.Login(ctx, subdomain, gofakeit.Email(), "@paSSw0rd", false)

			require.Error(t, err)
			require.ErrorIs(t, assert.AnError, err)
		})

	t.Run("should return error when userService.GetUserByOrgIDEmail returns error",
		func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			subdomain := gofakeit.LetterN(30)
			o := organization.Organization{ID: gofakeit.Int64()}
			email := gofakeit.Email()

			orgService := organization.NewMockService(t)
			orgService.On("GetOrganizationBySubdomain", ctx, subdomain).Return(o, nil)

			userService := user.NewMockService(t)
			userService.On("GetUserByOrgIDEmail", ctx, o.ID, email).Return(user.User{}, assert.AnError)

			authService := auth.NewService("secret", nil, orgService, userService, nil)
			_, _, err := authService.Login(ctx, subdomain, email, validPassword, false)

			require.Error(t, err)
			require.ErrorIs(t, assert.AnError, err)
		})

	t.Run("should return error when userService.GetUserByOrgIDEmail returns not found error",
		func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			subdomain := gofakeit.LetterN(30)
			o := organization.Organization{ID: gofakeit.Int64()}
			email := gofakeit.Email()

			orgService := organization.NewMockService(t)
			orgService.On("GetOrganizationBySubdomain", ctx, subdomain).Return(o, nil)

			userService := user.NewMockService(t)
			userService.On("GetUserByOrgIDEmail", ctx, o.ID, email).Return(user.User{}, base.NewNotFoundError("not found"))

			authService := auth.NewService("secret", nil, orgService, userService, nil)
			_, _, err := authService.Login(ctx, subdomain, email, validPassword, false)

			require.Error(t, err)
			require.ErrorIs(t, auth.ErrInvalidCredentials, err)
		})

	t.Run("should return error when user is disabled", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		subdomain := gofakeit.LetterN(30)
		email := gofakeit.Email()
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
		require.NoError(t, err)

		now := time.Now()
		u := user.User{
			ID:           gofakeit.Int64(),
			PasswordHash: string(passwordHash),
			DisabledAt:   &now,
		}
		o := organization.Organization{ID: gofakeit.Int64()}

		orgService := organization.NewMockService(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).Return(o, nil)

		userService := user.NewMockService(t)
		userService.On("GetUserByOrgIDEmail", ctx, o.ID, email).Return(u, nil)

		authService := auth.NewService("secret", nil, orgService, userService, nil)
		_, _, err = authService.Login(ctx, subdomain, email, validPassword, false)

		require.Error(t, err)
		require.ErrorIs(t, auth.ErrUserDisabled, err)
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		subdomain := gofakeit.LetterN(30)
		email := gofakeit.Email()
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
		require.NoError(t, err)

		u := user.User{
			ID:           gofakeit.Int64(),
			PasswordHash: string(passwordHash),
		}
		o := organization.Organization{ID: gofakeit.Int64()}

		orgService := organization.NewMockService(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).Return(o, nil)

		userService := user.NewMockService(t)
		userService.On("GetUserByOrgIDEmail", ctx, o.ID, email).Return(u, nil)

		authService := auth.NewService("secret", nil, orgService, userService, nil)
		_, _, err = authService.Login(ctx, subdomain, email, validPassword+"ZZZ", false)

		require.Error(t, err)
		require.ErrorIs(t, auth.ErrInvalidCredentials, err)
	})

	t.Run("should return error when sessionManager.CreateSession returns error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		subdomain := gofakeit.LetterN(30)
		email := gofakeit.Email()
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
		require.NoError(t, err)

		apiToken := gofakeit.UUID()
		u := user.User{
			ID:           gofakeit.Int64(),
			PasswordHash: string(passwordHash),
			APIToken:     &apiToken,
		}
		o := organization.Organization{ID: gofakeit.Int64()}

		orgService := organization.NewMockService(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).Return(o, nil)

		userService := user.NewMockService(t)
		userService.On("GetUserByOrgIDEmail", ctx, o.ID, email).Return(u, nil)

		sessionManager := session.NewMockSessionManager(t)
		sessionManager.On("CreateSession", ctx, u.ID, o.ID, fake.MockString, apiToken,
			auth.DefaultSessionTTL).Return(assert.AnError)

		authService := auth.NewService("jwt_secret", nil, orgService, userService, sessionManager)
		_, _, err = authService.Login(ctx, subdomain, email, validPassword, false)

		require.Error(t, err)
		require.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return the jwt token and default ttl", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		subdomain := gofakeit.LetterN(30)
		email := gofakeit.Email()
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
		require.NoError(t, err)

		apiToken := gofakeit.UUID()
		u := user.User{
			ID:           gofakeit.Int64(),
			PasswordHash: string(passwordHash),
			APIToken:     &apiToken,
		}
		o := organization.Organization{ID: gofakeit.Int64()}

		orgService := organization.NewMockService(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).Return(o, nil)

		userService := user.NewMockService(t)
		userService.On("GetUserByOrgIDEmail", ctx, o.ID, email).Return(u, nil)

		sessionManager := session.NewMockSessionManager(t)
		sessionManager.On("CreateSession", ctx, u.ID, o.ID, fake.MockString, apiToken,
			auth.DefaultSessionTTL).Return(nil)

		authService := auth.NewService("jwt_secret", nil, orgService, userService, sessionManager)
		token, ttl, err := authService.Login(ctx, subdomain, email, validPassword, false)

		require.NoError(t, err)
		require.NotEmpty(t, token)
		assert.Equal(t, auth.DefaultSessionTTL, ttl)
	})

	t.Run("should return the jwt token and remember-me ttl", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		subdomain := gofakeit.LetterN(30)
		email := gofakeit.Email()
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
		require.NoError(t, err)

		apiToken := gofakeit.UUID()
		u := user.User{
			ID:           gofakeit.Int64(),
			PasswordHash: string(passwordHash),
			APIToken:     &apiToken,
		}
		o := organization.Organization{ID: gofakeit.Int64()}

		orgService := organization.NewMockService(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).Return(o, nil)

		userService := user.NewMockService(t)
		userService.On("GetUserByOrgIDEmail", ctx, o.ID, email).Return(u, nil)

		sessionManager := session.NewMockSessionManager(t)
		sessionManager.On("CreateSession", ctx, u.ID, o.ID, fake.MockString, apiToken,
			auth.RememberMeSessionTTL).Return(nil)

		authService := auth.NewService("jwt_secret", nil, orgService, userService, sessionManager)
		token, ttl, err := authService.Login(ctx, subdomain, email, validPassword, true)

		require.NoError(t, err)
		require.NotEmpty(t, token)
		assert.Equal(t, auth.RememberMeSessionTTL, ttl)
	})
}

func TestService_Logout(t *testing.T) {
	t.Parallel()

	t.Run("should return error when sessionManager.DeleteSession returns error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()

		sessionManager := session.NewMockSessionManager(t)
		sessionManager.On("DeleteSession", ctx, userID, orgID).Return(assert.AnError)

		authService := auth.NewService("", nil, nil, nil, sessionManager)
		err := authService.Logout(ctx, userID, orgID)

		require.Error(t, err)
		require.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should logout successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()

		sessionManager := session.NewMockSessionManager(t)
		sessionManager.On("DeleteSession", ctx, userID, orgID).Return(nil)

		authService := auth.NewService("", nil, nil, nil, sessionManager)
		err := authService.Logout(ctx, userID, orgID)

		require.NoError(t, err)
	})
}
