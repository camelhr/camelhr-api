package organization_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/domains/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_GetOrganizationByID(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)

		mockRepo.On("GetOrganizationByID", context.Background(), int64(1)).
			Return(organization.Organization{}, assert.AnError)

		_, err := service.GetOrganizationByID(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return an error when the organization is not found", func(t *testing.T) {
		t.Parallel()

		var notFoundErr *base.NotFoundError

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)

		mockRepo.On("GetOrganizationByID", context.Background(), int64(1)).
			Return(organization.Organization{}, sql.ErrNoRows)

		_, err := service.GetOrganizationByID(context.Background(), int64(1))
		require.Error(t, err)
		assert.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should return the organization by ID", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)

		org := organization.Organization{
			ID:        1,
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockRepo.On("GetOrganizationByID", context.Background(), org.ID).
			Return(org, nil)

		result, err := service.GetOrganizationByID(context.Background(), org.ID)
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestService_GetOrganizationBySubdomain(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the subdomain is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		subdomain := "#invalid-subdomain"

		_, err := service.GetOrganizationBySubdomain(context.Background(), subdomain)
		require.Error(t, err)
		assert.ErrorContains(t, err, "subdomain can only contain alphanumeric characters")
	})

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgSubdomain := randomOrganizationSubdomain()

		mockRepo.On("GetOrganizationBySubdomain", context.Background(), orgSubdomain).
			Return(organization.Organization{}, assert.AnError)

		_, err := service.GetOrganizationBySubdomain(context.Background(), orgSubdomain)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return an error when the organization is not found", func(t *testing.T) {
		t.Parallel()

		var notFoundErr *base.NotFoundError

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgSubdomain := randomOrganizationSubdomain()

		mockRepo.On("GetOrganizationBySubdomain", context.Background(), orgSubdomain).
			Return(organization.Organization{}, sql.ErrNoRows)

		_, err := service.GetOrganizationBySubdomain(context.Background(), orgSubdomain)
		require.Error(t, err)
		assert.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should return the organization by subdomain", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgSubdomain := randomOrganizationSubdomain()

		org := organization.Organization{
			ID:        1,
			Subdomain: orgSubdomain,
		}

		mockRepo.On("GetOrganizationBySubdomain", context.Background(), orgSubdomain).
			Return(org, nil)

		result, err := service.GetOrganizationBySubdomain(context.Background(), orgSubdomain)
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestService_GetOrganizationByName(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the organization name is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgName := "ørg1-non-ascii"

		_, err := service.GetOrganizationByName(context.Background(), orgName)
		require.Error(t, err)
		assert.ErrorContains(t, err, "organization name can only contain ascii characters")
	})

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgName := randomOrganizationName()

		mockRepo.On("GetOrganizationByName", context.Background(), orgName).
			Return(organization.Organization{}, assert.AnError)

		_, err := service.GetOrganizationByName(context.Background(), orgName)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return an error when the organization is not found", func(t *testing.T) {
		t.Parallel()

		var notFoundErr *base.NotFoundError

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgName := randomOrganizationName()

		mockRepo.On("GetOrganizationByName", context.Background(), orgName).
			Return(organization.Organization{}, sql.ErrNoRows)

		_, err := service.GetOrganizationByName(context.Background(), orgName)
		require.Error(t, err)
		assert.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should return the organization by name", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)

		org := organization.Organization{
			ID:   1,
			Name: randomOrganizationName(),
		}

		mockRepo.On("GetOrganizationByName", context.Background(), org.Name).
			Return(org, nil)

		result, err := service.GetOrganizationByName(context.Background(), org.Name)
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestService_CreateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the subdomain is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		subdomain := "#invalid-subdomain"

		_, err := service.CreateOrganization(context.Background(), subdomain, randomOrganizationName())
		require.Error(t, err)
		assert.ErrorContains(t, err, "subdomain can only contain alphanumeric characters")
	})

	t.Run("should return an error when the organization name is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgName := "ørg1"

		_, err := service.CreateOrganization(context.Background(), randomOrganizationSubdomain(), orgName)
		require.Error(t, err)
		assert.ErrorContains(t, err, "organization name can only contain ascii characters")
	})

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)

		mockRepo.On("CreateOrganization", context.Background(), "sub1", "org1").
			Return(organization.Organization{}, assert.AnError)

		_, err := service.CreateOrganization(context.Background(), "sub1", "org1")
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return the created organization", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)

		org := organization.Organization{
			Subdomain: randomOrganizationSubdomain(),
			Name:      randomOrganizationName(),
		}

		mockRepo.On("CreateOrganization", context.Background(), org.Subdomain, org.Name).
			Return(org, nil)

		result, err := service.CreateOrganization(context.Background(), org.Subdomain, org.Name)
		require.NoError(t, err)
		assert.Equal(t, org, result)
	})
}

func TestService_UpdateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the organization name is invalid", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgID := gofakeit.Int64()
		newOrgName := "ørg1"

		err := service.UpdateOrganization(context.Background(), orgID, newOrgName)
		require.Error(t, err)
		assert.ErrorContains(t, err, "organization name can only contain ascii characters")
	})

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgID := gofakeit.Int64()
		newOrgName := randomOrganizationName()

		mockRepo.On("UpdateOrganization", context.Background(), orgID, newOrgName).
			Return(assert.AnError)

		err := service.UpdateOrganization(context.Background(), orgID, newOrgName)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is updated", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgID := gofakeit.Int64()
		newOrgName := randomOrganizationName()

		mockRepo.On("UpdateOrganization", context.Background(), orgID, newOrgName).
			Return(nil)

		err := service.UpdateOrganization(context.Background(), orgID, newOrgName)
		require.NoError(t, err)
	})
}

func TestService_DeleteOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgID := gofakeit.Int64()

		mockRepo.On("DeleteOrganization", context.Background(), orgID).
			Return(assert.AnError)

		err := service.DeleteOrganization(context.Background(), orgID)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return error when session manager fails to delete sessions", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		mockSessionManager := session.NewMockSessionManager(t)
		service := organization.NewService(mockRepo, mockSessionManager)
		orgID := gofakeit.Int64()

		mockRepo.On("DeleteOrganization", context.Background(), orgID).
			Return(nil)

		mockSessionManager.On("DeleteAllOrgSessions", context.Background(), orgID).
			Return(assert.AnError)

		err := service.DeleteOrganization(context.Background(), orgID)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is deleted", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		mockSessionManager := session.NewMockSessionManager(t)
		service := organization.NewService(mockRepo, mockSessionManager)
		orgID := gofakeit.Int64()

		mockRepo.On("DeleteOrganization", context.Background(), orgID).
			Return(nil)
		mockSessionManager.On("DeleteAllOrgSessions", context.Background(), orgID).
			Return(nil)

		err := service.DeleteOrganization(context.Background(), orgID)
		require.NoError(t, err)
	})
}

func TestService_SuspendOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgID := gofakeit.Int64()
		comment := "test suspend"

		mockRepo.On("SuspendOrganization", context.Background(), orgID, comment).
			Return(assert.AnError)

		err := service.SuspendOrganization(context.Background(), orgID, comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is suspended", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgID := gofakeit.Int64()
		comment := "test suspend"

		mockRepo.On("SuspendOrganization", context.Background(), orgID, comment).
			Return(nil)

		err := service.SuspendOrganization(context.Background(), orgID, comment)
		require.NoError(t, err)
	})
}

func TestService_UnsuspendOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgID := gofakeit.Int64()
		comment := "test unsuspend"

		mockRepo.On("UnsuspendOrganization", context.Background(), orgID, comment).
			Return(assert.AnError)

		err := service.UnsuspendOrganization(context.Background(), orgID, comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is unsuspended", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgID := gofakeit.Int64()
		comment := "test unsuspend"

		mockRepo.On("UnsuspendOrganization", context.Background(), orgID, comment).
			Return(nil)

		err := service.UnsuspendOrganization(context.Background(), orgID, comment)
		require.NoError(t, err)
	})
}

func TestService_DisableOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgID := gofakeit.Int64()
		comment := "test disable"

		mockRepo.On("DisableOrganization", context.Background(), orgID, comment).
			Return(assert.AnError)

		err := service.DisableOrganization(context.Background(), orgID, comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is disabled", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgID := gofakeit.Int64()
		comment := "test disable"

		mockRepo.On("DisableOrganization", context.Background(), orgID, comment).
			Return(nil)

		err := service.DisableOrganization(context.Background(), orgID, comment)
		require.NoError(t, err)
	})
}

func TestService_EnableOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the repository call fails", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgID := gofakeit.Int64()
		comment := "test enable"

		mockRepo.On("EnableOrganization", context.Background(), orgID, comment).
			Return(assert.AnError)

		err := service.EnableOrganization(context.Background(), orgID, comment)
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil when the organization is enabled", func(t *testing.T) {
		t.Parallel()

		mockRepo := organization.NewMockRepository(t)
		service := organization.NewService(mockRepo, nil)
		orgID := gofakeit.Int64()
		comment := "test enable"

		mockRepo.On("EnableOrganization", context.Background(), orgID, comment).
			Return(nil)

		err := service.EnableOrganization(context.Background(), orgID, comment)
		require.NoError(t, err)
	})
}
