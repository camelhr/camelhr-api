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
	"github.com/camelhr/camelhr-api/internal/domains/user"
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
		subdomain := gofakeit.Word()

		orgService := organization.NewServiceMock(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).
			Return(organization.Organization{ID: gofakeit.Int64(), Subdomain: subdomain}, nil)

		authService := auth.NewService(nil, orgService, nil, "")
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
			subdomain := gofakeit.Word()

			orgService := organization.NewServiceMock(t)
			orgService.On("GetOrganizationBySubdomain", ctx, subdomain).
				Return(organization.Organization{}, assert.AnError)

			authService := auth.NewService(nil, orgService, nil, "")
			err := authService.Register(ctx, email, validPassword, subdomain, orgName)

			require.Error(t, err)
			require.ErrorIs(t, assert.AnError, err)
		})

	t.Run("should return error when transactor.WithTx returns error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		email := gofakeit.Email()
		orgName := gofakeit.Company()
		subdomain := gofakeit.Word()
		notFoundErr := base.NewNotFoundError("not found")

		orgService := organization.NewServiceMock(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).
			Return(organization.Organization{}, notFoundErr)

		transactor := database.NewTransactorMock(t)
		transactor.On("WithTx", ctx, mock.Anything).Return(assert.AnError)

		authService := auth.NewService(transactor, orgService, nil, "")
		err := authService.Register(ctx, email, validPassword, subdomain, orgName)

		require.Error(t, err)
		require.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should not return error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		email := gofakeit.Email()
		orgName := gofakeit.Company()
		subdomain := gofakeit.Word()
		notFoundErr := base.NewNotFoundError("not found")

		orgService := organization.NewServiceMock(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).
			Return(organization.Organization{}, notFoundErr)

		transactor := database.NewTransactorMock(t)
		transactor.On("WithTx", ctx, mock.Anything).Return(nil)

		authService := auth.NewService(transactor, orgService, nil, "")
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
			subdomain := gofakeit.Word()

			orgService := organization.NewServiceMock(t)
			orgService.On("GetOrganizationBySubdomain", ctx, subdomain).
				Return(organization.Organization{}, assert.AnError)

			authService := auth.NewService(nil, orgService, nil, "secret")
			_, err := authService.Login(ctx, subdomain, gofakeit.Email(), "@paSSw0rd")

			require.Error(t, err)
			require.ErrorIs(t, assert.AnError, err)
		})

	t.Run("should return error when userService.GetUserByOrgIDEmail returns error",
		func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			subdomain := gofakeit.Word()
			o := organization.Organization{ID: gofakeit.Int64()}
			email := gofakeit.Email()

			orgService := organization.NewServiceMock(t)
			orgService.On("GetOrganizationBySubdomain", ctx, subdomain).
				Return(o, nil)

			userService := user.NewServiceMock(t)
			userService.On("GetUserByOrgIDEmail", ctx, o.ID, email).
				Return(user.User{}, assert.AnError)

			authService := auth.NewService(nil, orgService, userService, "secret")
			_, err := authService.Login(ctx, subdomain, email, validPassword)

			require.Error(t, err)
			require.ErrorIs(t, assert.AnError, err)
		})

	t.Run("should return error when userService.GetUserByOrgIDEmail returns not found error",
		func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			subdomain := gofakeit.Word()
			o := organization.Organization{ID: gofakeit.Int64()}
			email := gofakeit.Email()

			orgService := organization.NewServiceMock(t)
			orgService.On("GetOrganizationBySubdomain", ctx, subdomain).
				Return(o, nil)

			userService := user.NewServiceMock(t)
			userService.On("GetUserByOrgIDEmail", ctx, o.ID, email).
				Return(user.User{}, base.NewNotFoundError("not found"))

			authService := auth.NewService(nil, orgService, userService, "secret")
			_, err := authService.Login(ctx, subdomain, email, validPassword)

			require.Error(t, err)
			require.ErrorIs(t, auth.ErrInvalidCredentials, err)
		})

	t.Run("should return error when user is disabled", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		subdomain := gofakeit.Word()
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

		orgService := organization.NewServiceMock(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).
			Return(o, nil)

		userService := user.NewServiceMock(t)
		userService.On("GetUserByOrgIDEmail", ctx, o.ID, email).
			Return(u, nil)

		authService := auth.NewService(nil, orgService, userService, "secret")
		_, err = authService.Login(ctx, subdomain, email, validPassword)

		require.Error(t, err)
		require.ErrorIs(t, auth.ErrUserDisabled, err)
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		subdomain := gofakeit.Word()
		email := gofakeit.Email()
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
		require.NoError(t, err)

		u := user.User{
			ID:           gofakeit.Int64(),
			PasswordHash: string(passwordHash),
		}
		o := organization.Organization{ID: gofakeit.Int64()}

		orgService := organization.NewServiceMock(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).
			Return(o, nil)

		userService := user.NewServiceMock(t)
		userService.On("GetUserByOrgIDEmail", ctx, o.ID, email).
			Return(u, nil)

		authService := auth.NewService(nil, orgService, userService, "secret")
		_, err = authService.Login(ctx, subdomain, email, validPassword+"ZZZ")

		require.Error(t, err)
		require.ErrorIs(t, auth.ErrInvalidCredentials, err)
	})

	t.Run("should return the jwt token", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		subdomain := gofakeit.Word()
		email := gofakeit.Email()
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
		require.NoError(t, err)

		u := user.User{
			ID:           gofakeit.Int64(),
			PasswordHash: string(passwordHash),
		}
		o := organization.Organization{ID: gofakeit.Int64()}

		orgService := organization.NewServiceMock(t)
		orgService.On("GetOrganizationBySubdomain", ctx, subdomain).
			Return(o, nil)

		userService := user.NewServiceMock(t)
		userService.On("GetUserByOrgIDEmail", ctx, o.ID, email).
			Return(u, nil)

		authService := auth.NewService(nil, orgService, userService, "jwt_secret")
		token, err := authService.Login(ctx, subdomain, email, validPassword)

		require.NoError(t, err)
		require.NotEmpty(t, token)
	})
}