package organization_test

import (
	"context"
	"testing"

	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganizationService_GetOrganizationByID(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewOrganizationService(mockRepo)

		mockRepo.On("GetOrganizationByID", context.TODO(), int64(1)).
			Return(organization.Organization{}, assert.AnError)

		_, err := service.GetOrganizationByID(context.TODO(), int64(1))
		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return the organization by ID", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewOrganizationService(mockRepo)

		org := organization.Organization{
			ID:   1,
			Name: "org1",
		}

		mockRepo.On("GetOrganizationByID", context.TODO(), int64(1)).
			Return(org, nil)

		result, err := service.GetOrganizationByID(context.TODO(), int64(1))
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestOrganizationService_GetOrganizationByName(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewOrganizationService(mockRepo)

		mockRepo.On("GetOrganizationByName", context.TODO(), "org1").
			Return(organization.Organization{}, assert.AnError)

		_, err := service.GetOrganizationByName(context.TODO(), "org1")
		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return the organization by name", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewOrganizationService(mockRepo)

		org := organization.Organization{
			ID:   1,
			Name: "org1",
		}

		mockRepo.On("GetOrganizationByName", context.TODO(), "org1").
			Return(org, nil)

		result, err := service.GetOrganizationByName(context.TODO(), "org1")
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestOrganizationService_CreateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewOrganizationService(mockRepo)

		org := organization.Organization{
			Name: "org1",
		}

		mockRepo.On("CreateOrganization", context.TODO(), org).
			Return(int64(0), assert.AnError)

		_, err := service.CreateOrganization(context.TODO(), org)
		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return the ID of the created organization", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewOrganizationService(mockRepo)

		org := organization.Organization{
			Name: "org1",
		}

		mockRepo.On("CreateOrganization", context.TODO(), org).
			Return(int64(1), nil)

		result, err := service.CreateOrganization(context.TODO(), org)
		require.NoError(t, err)
		assert.Equal(t, int64(1), result)
	})
}

func TestOrganizationService_UpdateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewOrganizationService(mockRepo)

		org := organization.Organization{
			Name: "org1",
		}

		mockRepo.On("UpdateOrganization", context.TODO(), org).
			Return(assert.AnError)

		err := service.UpdateOrganization(context.TODO(), org)
		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is updated", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewOrganizationService(mockRepo)

		org := organization.Organization{
			Name: "org1",
		}

		mockRepo.On("UpdateOrganization", context.TODO(), org).
			Return(nil)

		err := service.UpdateOrganization(context.TODO(), org)
		require.NoError(t, err)
	})
}

func TestOrganizationService_DeleteOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewOrganizationService(mockRepo)

		mockRepo.On("DeleteOrganization", context.TODO(), int64(1)).
			Return(assert.AnError)

		err := service.DeleteOrganization(context.TODO(), int64(1))
		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is deleted", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewOrganizationService(mockRepo)

		mockRepo.On("DeleteOrganization", context.TODO(), int64(1)).
			Return(nil)

		err := service.DeleteOrganization(context.TODO(), int64(1))
		require.NoError(t, err)
	})
}
