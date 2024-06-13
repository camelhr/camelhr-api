package session_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/session"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionManager_CreateSession(t *testing.T) {
	t.Parallel()

	t.Run("should return error when org id is missing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := int64(0)
		jwt := gofakeit.UUID()
		apiToken := gofakeit.UUID()
		redisClient, _ := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		err := sessionManager.CreateSession(ctx, userID, orgID, jwt, apiToken, time.Hour)
		require.Error(t, err)
		require.ErrorIs(t, err, session.ErrMissingOrgID)
	})

	t.Run("should return error when user id is missing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := int64(0)
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()
		apiToken := gofakeit.UUID()
		redisClient, _ := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		err := sessionManager.CreateSession(ctx, userID, orgID, jwt, apiToken, time.Hour)
		require.Error(t, err)
		require.ErrorIs(t, err, session.ErrMissingUserID)
	})

	t.Run("should return error when both jwt & api-token is missing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := ""
		apiToken := ""
		redisClient, _ := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		err := sessionManager.CreateSession(ctx, userID, orgID, jwt, apiToken, time.Hour)
		require.Error(t, err)
		require.ErrorIs(t, err, session.ErrMissingToken)
	})

	t.Run("should return error when setting api token fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()
		apiToken := gofakeit.UUID()
		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectHSet("apiToken:"+apiToken, "org", orgID, "user", userID).SetErr(assert.AnError)

		err := sessionManager.CreateSession(ctx, userID, orgID, jwt, apiToken, time.Hour)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when setting session fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()
		apiToken := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectHSet("apiToken:"+apiToken, "org", orgID, "user", userID).SetVal(1)
		redisClientMock.ExpectExpire("apiToken:"+apiToken, time.Hour*24).SetVal(true)
		redisClientMock.ExpectHSet(fmt.Sprintf("session:org:%d:user:%d", orgID, userID), "org", orgID,
			"user", userID, "jwt", jwt, "apiToken", apiToken).SetErr(assert.AnError)

		err := sessionManager.CreateSession(ctx, userID, orgID, jwt, apiToken, time.Hour)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should create session successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()
		apiToken := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectHSet("apiToken:"+apiToken, "org", orgID, "user", userID).SetVal(1)
		redisClientMock.ExpectExpire("apiToken:"+apiToken, time.Hour*24).SetVal(true)
		redisClientMock.ExpectHSet(fmt.Sprintf("session:org:%d:user:%d", orgID, userID), "org", orgID,
			"user", userID, "jwt", jwt, "apiToken", apiToken).SetVal(1)
		redisClientMock.ExpectExpire(fmt.Sprintf("session:org:%d:user:%d", orgID, userID), time.Hour).SetVal(true)

		err := sessionManager.CreateSession(ctx, userID, orgID, jwt, apiToken, time.Hour)
		require.NoError(t, err)
	})

	t.Run("should create session without api token", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectHSet(fmt.Sprintf("session:org:%d:user:%d", orgID, userID), "org", orgID,
			"user", userID, "jwt", jwt).SetVal(1)
		redisClientMock.ExpectExpire(fmt.Sprintf("session:org:%d:user:%d", orgID, userID), time.Hour).SetVal(true)

		err := sessionManager.CreateSession(ctx, userID, orgID, jwt, "", time.Hour)
		require.NoError(t, err)
	})

	t.Run("should create session without jwt", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		apiToken := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectHSet("apiToken:"+apiToken, "org", orgID, "user", userID).SetVal(1)
		redisClientMock.ExpectExpire("apiToken:"+apiToken, time.Hour*24).SetVal(true)
		redisClientMock.ExpectHSet(fmt.Sprintf("session:org:%d:user:%d", orgID, userID), "org", orgID,
			"user", userID, "apiToken", apiToken).SetVal(1)
		redisClientMock.ExpectExpire(fmt.Sprintf("session:org:%d:user:%d", orgID, userID), time.Hour).SetVal(true)

		err := sessionManager.CreateSession(ctx, userID, orgID, "", apiToken, time.Hour)
		require.NoError(t, err)
	})
}

func TestSessionManager_ValidateJWTSession(t *testing.T) {
	t.Parallel()

	t.Run("should return error when org id is missing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := int64(0)
		jwt := gofakeit.UUID()
		redisClient, _ := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		err := sessionManager.ValidateJWTSession(ctx, userID, orgID, jwt)
		require.Error(t, err)
		require.ErrorIs(t, err, session.ErrMissingOrgID)
	})

	t.Run("should return error when user id is missing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := int64(0)
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()
		redisClient, _ := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		err := sessionManager.ValidateJWTSession(ctx, userID, orgID, jwt)
		require.Error(t, err)
		require.ErrorIs(t, err, session.ErrMissingUserID)
	})

	t.Run("should return error when jwt is missing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := ""
		redisClient, _ := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		err := sessionManager.ValidateJWTSession(ctx, userID, orgID, jwt)
		require.Error(t, err)
		require.ErrorIs(t, err, session.ErrMissingToken)
	})

	t.Run("should return error when session not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectExists(fmt.Sprintf("session:org:%d:user:%d", orgID, userID)).SetVal(0)

		err := sessionManager.ValidateJWTSession(ctx, userID, orgID, jwt)
		require.Error(t, err)
		require.ErrorIs(t, err, session.ErrInvalidSession)
	})

	t.Run("should return error when retrieving session jwt fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectExists(fmt.Sprintf("session:org:%d:user:%d", orgID, userID)).SetVal(1)
		redisClientMock.ExpectHGet(fmt.Sprintf("session:org:%d:user:%d", orgID, userID), "jwt").SetErr(assert.AnError)

		err := sessionManager.ValidateJWTSession(ctx, userID, orgID, jwt)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when session jwt mismatch", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectExists(fmt.Sprintf("session:org:%d:user:%d", orgID, userID)).SetVal(1)
		redisClientMock.ExpectHGet(fmt.Sprintf("session:org:%d:user:%d", orgID, userID), "jwt").SetVal("invalid-jwt")

		err := sessionManager.ValidateJWTSession(ctx, userID, orgID, jwt)
		require.Error(t, err)
		require.ErrorIs(t, err, session.ErrInvalidSession)
	})

	t.Run("should validate jwt session successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectExists(fmt.Sprintf("session:org:%d:user:%d", orgID, userID)).SetVal(1)
		redisClientMock.ExpectHGet(fmt.Sprintf("session:org:%d:user:%d", orgID, userID), "jwt").SetVal(jwt)

		err := sessionManager.ValidateJWTSession(ctx, userID, orgID, jwt)
		require.NoError(t, err)
	})
}

func TestSessionManager_ValidateAPITokenSession(t *testing.T) {
	t.Parallel()

	t.Run("should return error when api token is missing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		apiToken := ""
		redisClient, _ := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		_, _, err := sessionManager.ValidateAPITokenSession(ctx, apiToken)
		require.Error(t, err)
		require.ErrorIs(t, err, session.ErrMissingToken)
	})

	t.Run("should return error when user data retrieval for api token fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		apiToken := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectHGetAll("apiToken:" + apiToken).SetErr(assert.AnError)

		_, _, err := sessionManager.ValidateAPITokenSession(ctx, apiToken)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when user data not found for api token", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		apiToken := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectHGetAll("apiToken:" + apiToken).SetVal(map[string]string{})

		_, _, err := sessionManager.ValidateAPITokenSession(ctx, apiToken)
		require.Error(t, err)
		require.ErrorIs(t, err, session.ErrInvalidSession)
	})

	t.Run("should return error when session not found for user data of api-token", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		apiToken := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectHGetAll("apiToken:" + apiToken).SetVal(map[string]string{
			"org":  "1",
			"user": "1",
		})
		redisClientMock.ExpectExists("session:org:1:user:1").SetVal(0)

		_, _, err := sessionManager.ValidateAPITokenSession(ctx, apiToken)
		require.Error(t, err)
		require.ErrorIs(t, err, session.ErrInvalidSession)
	})

	t.Run("should return error when retrieving session api-token fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		apiToken := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectHGetAll("apiToken:" + apiToken).SetVal(map[string]string{
			"org":  "1",
			"user": "1",
		})
		redisClientMock.ExpectExists("session:org:1:user:1").SetVal(1)
		redisClientMock.ExpectHGet("session:org:1:user:1", "apiToken").SetErr(assert.AnError)

		_, _, err := sessionManager.ValidateAPITokenSession(ctx, apiToken)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when session api-token mismatch", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		apiToken := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectHGetAll("apiToken:" + apiToken).SetVal(map[string]string{
			"org":  "1",
			"user": "1",
		})
		redisClientMock.ExpectExists("session:org:1:user:1").SetVal(1)
		redisClientMock.ExpectHGet("session:org:1:user:1", "apiToken").SetVal("invalid-api-token")

		_, _, err := sessionManager.ValidateAPITokenSession(ctx, apiToken)
		require.Error(t, err)
		require.ErrorIs(t, err, session.ErrInvalidSession)
	})

	t.Run("should validate api-token session successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		apiToken := gofakeit.UUID()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectHGetAll("apiToken:" + apiToken).SetVal(map[string]string{
			"org":  "123",
			"user": "456",
		})
		redisClientMock.ExpectExists("session:org:123:user:456").SetVal(1)
		redisClientMock.ExpectHGet("session:org:123:user:456", "apiToken").SetVal(apiToken)

		userID, orgID, err := sessionManager.ValidateAPITokenSession(ctx, apiToken)
		require.NoError(t, err)
		require.Equal(t, int64(456), userID)
		require.Equal(t, int64(123), orgID)
	})
}

func TestSessionManager_DeleteSession(t *testing.T) {
	t.Parallel()

	t.Run("should return error when session not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectDel(fmt.Sprintf("session:org:%d:user:%d", orgID, userID)).SetErr(assert.AnError)

		err := sessionManager.DeleteSession(ctx, userID, orgID)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should delete session successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectDel(fmt.Sprintf("session:org:%d:user:%d", orgID, userID)).SetVal(1)

		err := sessionManager.DeleteSession(ctx, userID, orgID)
		require.NoError(t, err)
	})
}

func TestSessionManager_DeleteAllOrgSessions(t *testing.T) {
	t.Parallel()

	t.Run("should return error when retrieving session keys fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		orgID := gofakeit.Int64()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectKeys(fmt.Sprintf("session:org:%d:user:*", orgID)).SetErr(assert.AnError)

		err := sessionManager.DeleteAllOrgSessions(ctx, orgID)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error when deleting session keys fails", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		orgID := gofakeit.Int64()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectKeys(fmt.Sprintf("session:org:%d:user:*", orgID)).
			SetVal([]string{
				"session:org:1:user:1",
				"session:org:1:user:2",
			})
		redisClientMock.ExpectDel("session:org:1:user:1", "session:org:1:user:2").SetErr(assert.AnError)

		err := sessionManager.DeleteAllOrgSessions(ctx, orgID)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should delete all org sessions successfully", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		orgID := gofakeit.Int64()

		redisClient, redisClientMock := redismock.NewClientMock()
		sessionManager := session.NewRedisSessionManager(redisClient)

		redisClientMock.ExpectKeys(fmt.Sprintf("session:org:%d:user:*", orgID)).
			SetVal([]string{
				"session:org:1:user:1",
				"session:org:1:user:2",
			})
		redisClientMock.ExpectDel("session:org:1:user:1", "session:org:1:user:2").SetVal(1)

		err := sessionManager.DeleteAllOrgSessions(ctx, orgID)
		require.NoError(t, err)
	})
}
