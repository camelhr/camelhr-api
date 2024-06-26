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

		mockDB := database.NewMockDatabase(t)
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

		mockDB := database.NewMockDatabase(t)
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

		mockDB := database.NewMockDatabase(t)
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

		mockDB := database.NewMockDatabase(t)
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

		mockDB := database.NewMockDatabase(t)
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

		mockDB := database.NewMockDatabase(t)
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

		mockDB := database.NewMockDatabase(t)
		repo := organization.NewRepository(mockDB)
		org := organization.Organization{
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockDB.On("Exec", context.Background(), mock.Anything,
			tests.QueryMatcher("createOrganizationQuery"), org.Subdomain, org.Name).
			Return(assert.AnError)

		_, err := repo.CreateOrganization(context.Background(), org.Subdomain, org.Name)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return the created organization", func(t *testing.T) {
		t.Parallel()

		var emptyOrg organization.Organization

		mockDB := database.NewMockDatabase(t)
		repo := organization.NewRepository(mockDB)
		org := organization.Organization{
			ID:        gofakeit.Int64(),
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockDB.On("Exec", context.Background(), &emptyOrg,
			tests.QueryMatcher("createOrganizationQuery"), org.Subdomain, org.Name).
			Run(func(args mock.Arguments) {
				// populate the passed argument with the organization
				arg, ok := args.Get(1).(*organization.Organization)
				require.True(t, ok)
				*arg = org
			}).
			Return(nil)

		result, err := repo.CreateOrganization(context.Background(), org.Subdomain, org.Name)
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestRepository_UpdateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := organization.NewRepository(mockDB)
		orgID := gofakeit.Int64()
		orgName := randomOrganizationName()

		mockDB.On("Exec", context.Background(), nil,
			tests.QueryMatcher("updateOrganizationQuery"), orgID, orgName).
			Return(assert.AnError)

		err := repo.UpdateOrganization(context.Background(), orgID, orgName)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is updated", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
		repo := organization.NewRepository(mockDB)
		orgID := gofakeit.Int64()
		orgName := randomOrganizationName()

		mockDB.On("Exec", context.Background(), nil,
			tests.QueryMatcher("updateOrganizationQuery"), orgID, orgName).
			Return(nil)

		err := repo.UpdateOrganization(context.Background(), orgID, orgName)
		require.NoError(t, err)
	})
}

func TestRepository_DeleteOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mockDB := database.NewMockDatabase(t)
		repo := organization.NewRepository(mockDB)

		mockDB.On("Exec", ctx, nil, tests.QueryMatcher("deleteOrganizationQuery"), int64(1), "test delete").
			Return(assert.AnError)

		err := repo.DeleteOrganization(ctx, 1, "test delete")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is deleted", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mockDB := database.NewMockDatabase(t)
		repo := organization.NewRepository(mockDB)

		mockDB.On("Exec", ctx, nil, tests.QueryMatcher("deleteOrganizationQuery"), int64(1), "test delete").
			Return(nil)

		err := repo.DeleteOrganization(ctx, 1, "test delete")
		require.NoError(t, err)
	})
}

func TestRepository_SuspendOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewMockDatabase(t)
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

		mockDB := database.NewMockDatabase(t)
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

		mockDB := database.NewMockDatabase(t)
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

		mockDB := database.NewMockDatabase(t)
		repo := organization.NewRepository(mockDB)

		mockDB.On("Exec", context.Background(), nil,
			tests.QueryMatcher("unsuspendOrganizationQuery"), int64(1), "test unsuspended").
			Return(nil)

		err := repo.UnsuspendOrganization(context.Background(), 1, "test unsuspended")
		require.NoError(t, err)
	})
}

func randomOrganizationSubdomain() string {
	return gofakeit.LetterN(uint(gofakeit.Number(1, 30)))
}

func randomOrganizationName() string {
	return fmt.Sprint(gofakeit.LetterN(8), " ", gofakeit.Company())
}
