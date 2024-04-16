package organization_test

import (
	"context"
	"strings"
	"testing"

	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var queryMatcher = func(queryLabel string) any {
	return mock.MatchedBy(func(a any) bool {
		if query, ok := a.(string); ok {
			return strings.Contains(query, queryLabel)
		}

		return false
	})
}

func TestOrganizationRepository_ListOrganizations(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewOrganizationRepository(mockDB)

		mockDB.On("List", context.TODO(), mock.Anything, queryMatcher("listOrganizationsQuery")).
			Return(assert.AnError)

		_, err := repo.ListOrganizations(context.TODO())
		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return a list of organizations", func(t *testing.T) {
		t.Parallel()

		var emptyOrgList []organization.Organization

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewOrganizationRepository(mockDB)

		orgs := []organization.Organization{
			{
				ID:   1,
				Name: "org1",
			},
			{
				ID:   2,
				Name: "org2",
			},
		}

		mockDB.On("List", context.TODO(), &emptyOrgList, queryMatcher("listOrganizationsQuery")).
			Run(func(args mock.Arguments) {
				arg, ok := args.Get(1).(*[]organization.Organization)
				require.True(t, ok)
				*arg = orgs
			}).Return(nil)

		result, err := repo.ListOrganizations(context.TODO())
		require.NoError(t, err)
		assert.Equal(t, orgs, result)
	})
}

func TestOrganizationRepository_GetOrganizationByID(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewOrganizationRepository(mockDB)

		mockDB.On("Get", context.TODO(), mock.Anything, queryMatcher("getOrganizationByIDQuery"), int64(1)).
			Return(assert.AnError)

		_, err := repo.GetOrganizationByID(context.TODO(), 1)
		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return an organization", func(t *testing.T) {
		t.Parallel()

		var emptyOrg organization.Organization

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewOrganizationRepository(mockDB)

		org := organization.Organization{
			ID:   1,
			Name: "org1",
		}

		mockDB.On("Get", context.TODO(), &emptyOrg, queryMatcher("getOrganizationByIDQuery"), int64(1)).
			Run(func(args mock.Arguments) {
				arg, ok := args.Get(1).(*organization.Organization)
				require.True(t, ok)
				*arg = org
			}).Return(nil)

		result, err := repo.GetOrganizationByID(context.TODO(), 1)
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestOrganizationRepository_GetOrganizationByName(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewOrganizationRepository(mockDB)

		mockDB.On("Get", context.TODO(), mock.Anything, queryMatcher("getOrganizationByNameQuery"), "org1").
			Return(assert.AnError)

		_, err := repo.GetOrganizationByName(context.TODO(), "org1")
		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return an organization", func(t *testing.T) {
		t.Parallel()

		var emptyOrg organization.Organization

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewOrganizationRepository(mockDB)

		org := organization.Organization{
			ID:   1,
			Name: "org1",
		}

		mockDB.On("Get", context.TODO(), &emptyOrg, queryMatcher("getOrganizationByNameQuery"), "org1").
			Run(func(args mock.Arguments) {
				arg, ok := args.Get(1).(*organization.Organization)
				require.True(t, ok)
				*arg = org
			}).Return(nil)

		result, err := repo.GetOrganizationByName(context.TODO(), "org1")
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestOrganizationRepository_CreateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewOrganizationRepository(mockDB)

		mockDB.On("Exec", context.TODO(), mock.Anything, queryMatcher("createOrganizationQuery"), "org1").
			Return(assert.AnError)

		_, err := repo.CreateOrganization(context.TODO(), organization.Organization{Name: "org1"})
		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return the organization ID", func(t *testing.T) {
		t.Parallel()

		var id int64

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewOrganizationRepository(mockDB)

		mockDB.On("Exec", context.TODO(), &id, queryMatcher("createOrganizationQuery"), "org1").
			Return(nil)

		result, err := repo.CreateOrganization(context.TODO(), organization.Organization{Name: "org1"})
		require.NoError(t, err)
		assert.Equal(t, id, result)
	})
}

func TestOrganizationRepository_UpdateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewOrganizationRepository(mockDB)

		mockDB.On("Exec", context.TODO(), nil, queryMatcher("updateOrganizationQuery"), "org1").
			Return(assert.AnError)

		err := repo.UpdateOrganization(context.TODO(), organization.Organization{Name: "org1"})
		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is updated", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewOrganizationRepository(mockDB)

		mockDB.On("Exec", context.TODO(), nil, queryMatcher("updateOrganizationQuery"), "org1").
			Return(nil)

		err := repo.UpdateOrganization(context.TODO(), organization.Organization{Name: "org1"})
		require.NoError(t, err)
	})
}

func TestOrganizationRepository_DeleteOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the database call fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &database.DatabaseMock{}
		repo := organization.NewOrganizationRepository(mockDB)

		mockDB.On("Exec", context.TODO(), nil, queryMatcher("deleteOrganizationQuery"), int64(1)).
			Return(assert.AnError)

		err := repo.DeleteOrganization(context.TODO(), 1)
		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is deleted", func(t *testing.T) {
		t.Parallel()

		mockDB := database.NewDatabaseMock(t)
		repo := organization.NewOrganizationRepository(mockDB)

		mockDB.On("Exec", context.TODO(), nil, queryMatcher("deleteOrganizationQuery"), int64(1)).
			Return(nil)

		err := repo.DeleteOrganization(context.TODO(), 1)
		require.NoError(t, err)
	})
}
