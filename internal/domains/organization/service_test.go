package organization_test

import (
	"context"
	"testing"

	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_GetOrganizationByID(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		mockRepo.On("GetOrganizationByID", context.TODO(), int64(1)).
			Return(organization.Organization{}, assert.AnError)

		_, err := service.GetOrganizationByID(context.TODO(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return the organization by ID", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		org := organization.Organization{
			ID:        1,
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockRepo.On("GetOrganizationByID", context.TODO(), int64(1)).
			Return(org, nil)

		result, err := service.GetOrganizationByID(context.TODO(), int64(1))
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestService_GetOrganizationBySubdomain(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)
		orgSubdomain := randomOrganizationSubdomain()

		mockRepo.On("GetOrganizationBySubdomain", context.TODO(), orgSubdomain).
			Return(organization.Organization{}, assert.AnError)

		_, err := service.GetOrganizationBySubdomain(context.TODO(), orgSubdomain)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return the organization by subdomain", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)
		orgSubdomain := randomOrganizationSubdomain()

		org := organization.Organization{
			ID:        1,
			Subdomain: orgSubdomain,
		}

		mockRepo.On("GetOrganizationBySubdomain", context.TODO(), orgSubdomain).
			Return(org, nil)

		result, err := service.GetOrganizationBySubdomain(context.TODO(), orgSubdomain)
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestService_GetOrganizationByName(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)
		orgName := randomOrganizationName()

		mockRepo.On("GetOrganizationByName", context.TODO(), orgName).
			Return(organization.Organization{}, assert.AnError)

		_, err := service.GetOrganizationByName(context.TODO(), orgName)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return the organization by name", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		org := organization.Organization{
			ID:   1,
			Name: randomOrganizationName(),
		}

		mockRepo.On("GetOrganizationByName", context.TODO(), org.Name).
			Return(org, nil)

		result, err := service.GetOrganizationByName(context.TODO(), org.Name)
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestService_CreateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		org := organization.Organization{
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockRepo.On("CreateOrganization", context.TODO(), org).
			Return(int64(0), assert.AnError)

		_, err := service.CreateOrganization(context.TODO(), org)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return the ID of the created organization", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		org := organization.Organization{
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockRepo.On("CreateOrganization", context.TODO(), org).
			Return(int64(1), nil)

		result, err := service.CreateOrganization(context.TODO(), org)
		require.NoError(t, err)
		assert.Equal(t, int64(1), result)
	})
}

func TestService_UpdateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		org := organization.Organization{
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockRepo.On("UpdateOrganization", context.TODO(), org).
			Return(assert.AnError)

		err := service.UpdateOrganization(context.TODO(), org)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is updated", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		org := organization.Organization{
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockRepo.On("UpdateOrganization", context.TODO(), org).
			Return(nil)

		err := service.UpdateOrganization(context.TODO(), org)
		require.NoError(t, err)
	})
}

func TestService_DeleteOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		mockRepo.On("DeleteOrganization", context.TODO(), int64(1)).
			Return(assert.AnError)

		err := service.DeleteOrganization(context.TODO(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is deleted", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		mockRepo.On("DeleteOrganization", context.TODO(), int64(1)).
			Return(nil)

		err := service.DeleteOrganization(context.TODO(), int64(1))
		require.NoError(t, err)
	})
}

func TestService_SuspendOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		mockRepo.On("SuspendOrganization", context.TODO(), int64(1), "test suspend").
			Return(assert.AnError)

		err := service.SuspendOrganization(context.TODO(), int64(1), "test suspend")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is suspended", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		mockRepo.On("SuspendOrganization", context.TODO(), int64(1), "test suspend").
			Return(nil)

		err := service.SuspendOrganization(context.TODO(), int64(1), "test suspend")
		require.NoError(t, err)
	})
}

func TestService_UnsuspendOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		mockRepo.On("UnsuspendOrganization", context.TODO(), int64(1), "test unsuspend").
			Return(assert.AnError)

		err := service.UnsuspendOrganization(context.TODO(), int64(1), "test unsuspend")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is unsuspended", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		mockRepo.On("UnsuspendOrganization", context.TODO(), int64(1), "test unsuspend").
			Return(nil)

		err := service.UnsuspendOrganization(context.TODO(), int64(1), "test unsuspend")
		require.NoError(t, err)
	})
}

func TestService_BlacklistOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		mockRepo.On("BlacklistOrganization", context.TODO(), int64(1), "test blacklist").
			Return(assert.AnError)

		err := service.BlacklistOrganization(context.TODO(), int64(1), "test blacklist")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is blacklisted", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		mockRepo.On("BlacklistOrganization", context.TODO(), int64(1), "test blacklist").
			Return(nil)

		err := service.BlacklistOrganization(context.TODO(), int64(1), "test blacklist")
		require.NoError(t, err)
	})
}

func TestService_UnblacklistOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		mockRepo.On("UnblacklistOrganization", context.TODO(), int64(1), "test unblacklist").
			Return(assert.AnError)

		err := service.UnblacklistOrganization(context.TODO(), int64(1), "test unblacklist")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is unblacklisted", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewRepositoryMock(t)
		service := organization.NewService(mockRepo)

		mockRepo.On("UnblacklistOrganization", context.TODO(), int64(1), "test unblacklist").
			Return(nil)

		err := service.UnblacklistOrganization(context.TODO(), int64(1), "test unblacklist")
		require.NoError(t, err)
	})
}
