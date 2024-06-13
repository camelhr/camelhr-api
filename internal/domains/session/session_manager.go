package session

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	sessionHKeyFormat      = "session:org:%v:user:%v"
	apiTokenHKeyFormat     = "apiToken:%s"
	apiTokenHKeyDefaultTTL = time.Hour * 24

	orgKey      = "org"
	userKey     = "user"
	jwtKey      = "jwt"
	apiTokenKey = "apiToken"
)

var (
	ErrInvalidSession = errors.New("invalid session")
	ErrMissingUserID  = errors.New("missing user id")
	ErrMissingOrgID   = errors.New("missing org id")
	ErrMissingToken   = errors.New("missing token")
)

// SessionManager is an interface for managing user sessions.
// It provides methods to create, validate and delete sessions.
// The session can be validated using either JWT or API token.
type SessionManager interface {
	// CreateSession creates a new session for the user of the given organization
	// The session can be created using either JWT or API token
	CreateSession(ctx context.Context, userID, orgID int64, jwt, apiToken string, ttl time.Duration) error

	// ValidateJWTSession validates the JWT session for the user of the given organization
	ValidateJWTSession(ctx context.Context, userID, orgID int64, jwt string) error

	// ValidateAPITokenSession validates the API token session and returns the associated userID and orgID
	ValidateAPITokenSession(ctx context.Context, apiToken string) (int64, int64, error)

	// DeleteSession deletes the session for the user of the given organization
	DeleteSession(ctx context.Context, userID, orgID int64) error

	// DeleteAllOrgSessions deletes all user sessions under the given organization
	DeleteAllOrgSessions(ctx context.Context, orgID int64) error
}

type sessionManager struct {
	redisClient *redis.Client
}

func NewRedisSessionManager(redisClient *redis.Client) SessionManager {
	return &sessionManager{redisClient}
}

func (m *sessionManager) CreateSession(
	ctx context.Context,
	userID, orgID int64,
	jwt, apiToken string,
	ttl time.Duration,
) error {
	if err := m.validateParams(userID, orgID, jwt, apiToken); err != nil {
		return err
	}

	sessionHKey := fmt.Sprintf(sessionHKeyFormat, orgID, userID)
	keyValues := []any{orgKey, orgID, userKey, userID}

	if jwt != "" {
		keyValues = append(keyValues, jwtKey, jwt)
	}

	if apiToken != "" {
		keyValues = append(keyValues, apiTokenKey, apiToken)

		// store user data for apiToken
		apiTokenHKey := fmt.Sprintf(apiTokenHKeyFormat, apiToken)
		if err := m.redisClient.HSet(ctx, apiTokenHKey, orgKey, orgID, userKey, userID).Err(); err != nil {
			return fmt.Errorf("failed to set api-token hash key for user:%d org:%d: %w", userID, orgID, err)
		}

		// set expiry for api-token hash key
		if err := m.redisClient.Expire(ctx, apiTokenHKey, apiTokenHKeyDefaultTTL).Err(); err != nil {
			return fmt.Errorf("failed to set api-token hash key expiry for user:%d org:%d: %w", userID, orgID, err)
		}
	}

	// store user session
	if err := m.redisClient.HSet(ctx, sessionHKey, keyValues...).Err(); err != nil {
		return fmt.Errorf("failed to persist session for user:%d org:%d: %w", userID, orgID, err)
	}

	// set expiry for session hash key
	if err := m.redisClient.Expire(ctx, sessionHKey, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set session hash key expiry for user:%d org:%d: %w", userID, orgID, err)
	}

	return nil
}

func (m *sessionManager) ValidateJWTSession(ctx context.Context, userID, orgID int64, jwt string) error {
	if err := m.validateParams(userID, orgID, jwt, ""); err != nil {
		return err
	}

	sessionHKey := fmt.Sprintf(sessionHKeyFormat, orgID, userID)

	if m.redisClient.Exists(ctx, sessionHKey).Val() == 0 {
		return fmt.Errorf("session not found for user:%d org:%d: %w",
			userID, orgID, ErrInvalidSession)
	}

	sessionJWT, err := m.redisClient.HGet(ctx, sessionHKey, jwtKey).Result()
	if err != nil {
		return fmt.Errorf("failed to retrieve session jwt for user:%d org:%d: %w",
			userID, orgID, err)
	}

	if sessionJWT != jwt {
		return fmt.Errorf("session jwt mismatch for user:%d org:%d: %w",
			userID, orgID, ErrInvalidSession)
	}

	return nil
}

func (m *sessionManager) ValidateAPITokenSession(ctx context.Context, apiToken string) (int64, int64, error) {
	if apiToken == "" {
		return 0, 0, ErrMissingToken
	}

	apiTokenHKey := fmt.Sprintf(apiTokenHKeyFormat, apiToken)

	userData, err := m.redisClient.HGetAll(ctx, apiTokenHKey).Result()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to retrieve user data of api-token for user:%s org:%s: %w",
			userData[userKey], userData[orgKey], err)
	}

	if len(userData) == 0 {
		return 0, 0, fmt.Errorf("user data not found for api-token: %w", ErrInvalidSession)
	}

	sessionHKey := fmt.Sprintf(sessionHKeyFormat, userData[orgKey], userData[userKey])
	if m.redisClient.Exists(ctx, sessionHKey).Val() == 0 {
		return 0, 0, fmt.Errorf("session not found for user:%s org:%s: %w",
			userData[userKey], userData[orgKey], ErrInvalidSession)
	}

	sessionAPIToken, err := m.redisClient.HGet(ctx, sessionHKey, apiTokenKey).Result()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to retrieve session api-token for user:%s org:%s: %w",
			userData[userKey], userData[orgKey], err)
	}

	if sessionAPIToken != apiToken {
		return 0, 0, fmt.Errorf("session api-token mismatch for user:%s org:%s: %w",
			userData[userKey], userData[orgKey], ErrInvalidSession)
	}

	userID, _ := strconv.ParseInt(userData[userKey], 10, 64)
	orgID, _ := strconv.ParseInt(userData[orgKey], 10, 64)

	return userID, orgID, nil
}

func (m *sessionManager) DeleteSession(ctx context.Context, userID, orgID int64) error {
	sessionHKey := fmt.Sprintf(sessionHKeyFormat, orgID, userID)
	if err := m.redisClient.Del(ctx, sessionHKey).Err(); err != nil {
		return fmt.Errorf("failed to delete session for user: %d org: %d: %w",
			userID, orgID, err)
	}

	return nil
}

func (m *sessionManager) DeleteAllOrgSessions(ctx context.Context, orgID int64) error {
	sessionHKeys, err := m.redisClient.Keys(ctx, fmt.Sprintf(sessionHKeyFormat, orgID, "*")).Result()
	if err != nil {
		return fmt.Errorf("failed to retrieve all session keys for org: %d: %w",
			orgID, err)
	}

	if err := m.redisClient.Del(ctx, sessionHKeys...).Err(); err != nil {
		return fmt.Errorf("failed to delete session keys for org: %d: %w",
			orgID, err)
	}

	return nil
}

func (m *sessionManager) validateParams(userID, orgID int64, jwt, apiToken string) error {
	if userID == 0 {
		return ErrMissingUserID
	}

	if orgID == 0 {
		return ErrMissingOrgID
	}

	// either jwt or apiToken should be present
	if jwt == "" && apiToken == "" {
		return ErrMissingToken
	}

	return nil
}
