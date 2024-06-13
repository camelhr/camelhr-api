package session_test

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/session"
)

func (s *SessionTestSuite) TestSessionManagerIntegration_CreateSession() {
	s.Run("should create a new session for the user of the given organization", func() {
		s.T().Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()
		apiToken := gofakeit.UUID()
		sessionKey := fmt.Sprintf("session:org:%v:user:%v", orgID, userID)
		sessionManager := session.NewRedisSessionManager(s.RedisClient)

		err := sessionManager.CreateSession(ctx, userID, orgID, jwt, apiToken, time.Hour)
		s.Require().NoError(err)

		sessionData := s.RedisClient.HGetAll(ctx, sessionKey).Val()
		s.Len(sessionData, 4)
		s.Equal(jwt, sessionData["jwt"])
		s.Equal(apiToken, sessionData["apiToken"])
		s.Equal(strconv.FormatInt(userID, 10), sessionData["user"])
		s.Equal(strconv.FormatInt(orgID, 10), sessionData["org"])

		ttl := s.RedisClient.TTL(ctx, sessionKey).Val()
		s.Equal(time.Hour, ttl)
	})

	s.Run("should overwrite the existing session for new toke", func() {
		s.T().Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()
		apiToken := gofakeit.UUID()
		sessionKey := fmt.Sprintf("session:org:%v:user:%v", orgID, userID)
		sessionManager := session.NewRedisSessionManager(s.RedisClient)

		err := sessionManager.CreateSession(ctx, userID, orgID, jwt, apiToken, time.Hour)
		s.Require().NoError(err)

		newJWT := gofakeit.UUID()
		newAPIToken := gofakeit.UUID()
		err = sessionManager.CreateSession(ctx, userID, orgID, newJWT, newAPIToken, time.Hour)
		s.Require().NoError(err)

		sessionData := s.RedisClient.HGetAll(ctx, sessionKey).Val()
		s.Len(sessionData, 4)
		s.Equal(newJWT, sessionData["jwt"])
		s.Equal(newAPIToken, sessionData["apiToken"])
		s.Equal(strconv.FormatInt(userID, 10), sessionData["user"])
		s.Equal(strconv.FormatInt(orgID, 10), sessionData["org"])

		ttl := s.RedisClient.TTL(ctx, sessionKey).Val()
		s.Equal(time.Hour, ttl)
	})
}

func (s *SessionTestSuite) TestSessionManagerIntegration_ValidateJWTSession() {
	s.Run("should validate the JWT session for the user of the given organization", func() {
		s.T().Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()
		apiToken := gofakeit.UUID()
		sessionKey := fmt.Sprintf("session:org:%v:user:%v", orgID, userID)
		sessionManager := session.NewRedisSessionManager(s.RedisClient)

		err := sessionManager.CreateSession(ctx, userID, orgID, jwt, apiToken, time.Hour)
		s.Require().NoError(err)
		sessionData := s.RedisClient.HGetAll(ctx, sessionKey).Val()
		s.Len(sessionData, 4)
		s.Require().Equal(jwt, sessionData["jwt"])

		err = sessionManager.ValidateJWTSession(ctx, userID, orgID, jwt)
		s.Require().NoError(err)
	})
}

func (s *SessionTestSuite) TestSessionManagerIntegration_ValidateAPITokenSession() {
	s.Run("should validate the API token session", func() {
		s.T().Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()
		apiToken := gofakeit.UUID()
		sessionKey := fmt.Sprintf("session:org:%v:user:%v", orgID, userID)
		sessionManager := session.NewRedisSessionManager(s.RedisClient)

		err := sessionManager.CreateSession(ctx, userID, orgID, jwt, apiToken, time.Hour)
		s.Require().NoError(err)
		sessionData := s.RedisClient.HGetAll(ctx, sessionKey).Val()
		s.Len(sessionData, 4)
		s.Require().Equal(apiToken, sessionData["apiToken"])

		resultUserID, resultOrgID, err := sessionManager.ValidateAPITokenSession(ctx, apiToken)
		s.Require().NoError(err)
		s.Equal(userID, resultUserID)
		s.Equal(orgID, resultOrgID)
	})
}

func (s *SessionTestSuite) TestSessionManagerIntegration_DeleteSession() {
	s.Run("should delete the session for the user of the given organization", func() {
		s.T().Parallel()

		ctx := context.Background()
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		jwt := gofakeit.UUID()
		apiToken := gofakeit.UUID()
		sessionKey := fmt.Sprintf("session:org:%v:user:%v", orgID, userID)
		sessionManager := session.NewRedisSessionManager(s.RedisClient)

		err := sessionManager.CreateSession(ctx, userID, orgID, jwt, apiToken, time.Hour)
		s.Require().NoError(err)
		sessionData := s.RedisClient.HGetAll(ctx, sessionKey).Val()
		s.Len(sessionData, 4)

		err = sessionManager.DeleteSession(ctx, userID, orgID)
		s.Require().NoError(err)
		exist := s.RedisClient.Exists(ctx, sessionKey).Val()
		s.Zero(exist)
	})
}

func (s *SessionTestSuite) TestSessionManagerIntegration_DeleteAllOrgSessions() {
	s.Run("should delete all user sessions under the given organization", func() {
		s.T().Parallel()

		ctx := context.Background()
		orgID := gofakeit.Int64()
		sessionManager := session.NewRedisSessionManager(s.RedisClient)

		err := sessionManager.CreateSession(ctx, gofakeit.Int64(), orgID, gofakeit.UUID(), gofakeit.UUID(), time.Hour)
		s.Require().NoError(err)
		err = sessionManager.CreateSession(ctx, gofakeit.Int64(), orgID, gofakeit.UUID(), gofakeit.UUID(), time.Hour)
		s.Require().NoError(err)
		keys, err := s.RedisClient.Keys(ctx, fmt.Sprintf("session:org:%v:user:*", orgID)).Result()
		s.Require().NoError(err)
		s.Len(keys, 2)

		err = sessionManager.DeleteAllOrgSessions(ctx, orgID)
		s.Require().NoError(err)

		keys, err = s.RedisClient.Keys(ctx, fmt.Sprintf("session:org:%v:user:*", orgID)).Result()
		s.Require().NoError(err)
		s.Empty(keys)
	})
}
