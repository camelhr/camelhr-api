package organization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetOrganizationByID(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewRepository(mockDB)

		mockDB.On("Get", context.Background(), mock.Anything, tests.QueryMatcher("getOrganizationByIDQuery"), int64(1)).
			Return(assert.AnError)

		_, err := repo.GetOrganizationByID(context.Background(), 1)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return an organization", func(t *testing.T) {
		t.Parallel()

		var emptyOrg organization.Organization

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewRepository(mockDB)

		org := organization.Organization{
			ID:        1,
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockDB.On("Get", context.Background(), &emptyOrg, tests.QueryMatcher("getOrganizationByIDQuery"), int64(1)).
			Run(func(args mock.Arguments) {
				// populate the passed argument with the organization
				arg, ok := args.Get(1).(*organization.Organization)
				require.True(t, ok)
				*arg = org
			}).Return(nil)

		result, err := repo.GetOrganizationByID(context.Background(), 1)
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestRepository_GetOrganizationBySubdomain(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewRepository(mockDB)

		mockDB.On("Get", context.Background(), mock.Anything, tests.QueryMatcher("getOrganizationBySubdomainQuery"), "org1").
			Return(assert.AnError)

		_, err := repo.GetOrganizationBySubdomain(context.Background(), "org1")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return an organization", func(t *testing.T) {
		t.Parallel()

		var emptyOrg organization.Organization

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewRepository(mockDB)

		org := organization.Organization{
			ID:        1,
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockDB.On("Get", context.Background(), &emptyOrg, tests.QueryMatcher("getOrganizationBySubdomainQuery"), "org1").
			Run(func(args mock.Arguments) {
				// populate the passed argument with the organization
				arg, ok := args.Get(1).(*organization.Organization)
				require.True(t, ok)
				*arg = org
			}).Return(nil)

		result, err := repo.GetOrganizationBySubdomain(context.Background(), "org1")
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestRepository_GetOrganizationByName(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewRepository(mockDB)

		mockDB.On("Get", context.Background(), mock.Anything, tests.QueryMatcher("getOrganizationByNameQuery"), "org1").
			Return(assert.AnError)

		_, err := repo.GetOrganizationByName(context.Background(), "org1")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return an organization", func(t *testing.T) {
		t.Parallel()

		var emptyOrg organization.Organization

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewRepository(mockDB)

		org := organization.Organization{
			ID:        1,
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockDB.On("Get", context.Background(), &emptyOrg, tests.QueryMatcher("getOrganizationByNameQuery"), "org1").
			Run(func(args mock.Arguments) {
				// populate the passed argument with the organization
				arg, ok := args.Get(1).(*organization.Organization)
				require.True(t, ok)
				*arg = org
			}).Return(nil)

		result, err := repo.GetOrganizationByName(context.Background(), "org1")
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestRepository_CreateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewRepository(mockDB)
		org := organization.Organization{
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockDB.On("Exec", context.Background(), mock.Anything,
			tests.QueryMatcher("createOrganizationQuery"), org.Subdomain, org.Name).
			Return(assert.AnError)

		_, err := repo.CreateOrganization(context.Background(), org)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return the organization ID", func(t *testing.T) {
		t.Parallel()

		var id int64

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewRepository(mockDB)
		org := organization.Organization{
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockDB.On("Exec", context.Background(), &id,
			tests.QueryMatcher("createOrganizationQuery"), org.Subdomain, org.Name).
			Return(nil)

		result, err := repo.CreateOrganization(context.Background(), org)
		require.NoError(t, err)
		assert.Equal(t, id, result)
	})
}

func TestRepository_UpdateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewRepository(mockDB)
		org := organization.Organization{
			ID:        gofakeit.Int64(),
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockDB.On("Exec", context.Background(), nil,
			tests.QueryMatcher("updateOrganizationQuery"), org.ID, org.Subdomain, org.Name).
			Return(assert.AnError)

		err := repo.UpdateOrganization(context.Background(), org)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is updated", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewRepository(mockDB)
		org := organization.Organization{
			ID:        gofakeit.Int64(),
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockDB.On("Exec", context.Background(), nil,
			tests.QueryMatcher("updateOrganizationQuery"), org.ID, org.Subdomain, org.Name).
			Return(nil)

		err := repo.UpdateOrganization(context.Background(), org)
		require.NoError(t, err)
	})
}

func TestRepository_DeleteOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("deleteOrganizationQuery"), int64(1)).
			Return(assert.AnError)

		err := repo.DeleteOrganization(context.Background(), 1)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is deleted", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil, tests.QueryMatcher("deleteOrganizationQuery"), int64(1)).
			Return(nil)

		err := repo.DeleteOrganization(context.Background(), 1)
		require.NoError(t, err)
	})
}

func TestRepository_SuspendOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil,
			tests.QueryMatcher("suspendOrganizationQuery"), int64(1), "test suspended").
			Return(assert.AnError)

		err := repo.SuspendOrganization(context.Background(), 1, "test suspended")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is suspended", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil,
			tests.QueryMatcher("suspendOrganizationQuery"), int64(1), "test suspended").
			Return(nil)

		err := repo.SuspendOrganization(context.Background(), 1, "test suspended")
		require.NoError(t, err)
	})
}

func TestRepository_UnsuspendOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil,
			tests.QueryMatcher("unsuspendOrganizationQuery"), int64(1), "test unsuspended").
			Return(assert.AnError)

		err := repo.UnsuspendOrganization(context.Background(), 1, "test unsuspended")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is unsuspended", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil,
			tests.QueryMatcher("unsuspendOrganizationQuery"), int64(1), "test unsuspended").
			Return(nil)

		err := repo.UnsuspendOrganization(context.Background(), 1, "test unsuspended")
		require.NoError(t, err)
	})
}

func TestRepository_BlacklistOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil,
			tests.QueryMatcher("blacklistOrganizationQuery"), int64(1), "test blacklisted").
			Return(assert.AnError)

		err := repo.BlacklistOrganization(context.Background(), 1, "test blacklisted")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is blacklisted", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil,
			tests.QueryMatcher("blacklistOrganizationQuery"), int64(1), "test blacklisted").
			Return(nil)

		err := repo.BlacklistOrganization(context.Background(), 1, "test blacklisted")
		require.NoError(t, err)
	})
}

func TestRepository_UnblacklistOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil,
			tests.QueryMatcher("unblacklistOrganizationQuery"), int64(1), "test unblacklisted").
			Return(assert.AnError)

		err := repo.UnblacklistOrganization(context.Background(), 1, "test unblacklisted")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is unblacklisted", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil,
			tests.QueryMatcher("unblacklistOrganizationQuery"), int64(1), "test unblacklisted").
			Return(nil)

		err := repo.UnblacklistOrganization(context.Background(), 1, "test unblacklisted")
		require.NoError(t, err)
	})
}

func randomOrganizationSubdomain() string {
	return gofakeit.LetterN(uint(gofakeit.Number(1, 30)))
}

func randomOrganizationName() string {
	return fmt.Sprint(gofakeit.LetterN(8), " ", gofakeit.Company())
}
