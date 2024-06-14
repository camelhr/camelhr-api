package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AppClaims represents the claims in the jwt token.
type AppClaims struct {
	UserID       int64  `json:"user_id"`
	OrgID        int64  `json:"org_id"`
	OrgSubdomain string `json:"org_subdomain"`
	jwt.RegisteredClaims
}

// Validate validates the claims.
// It will be called by the jwt.ParseWithClaims after parsing the token.
func (c *AppClaims) Validate() error {
	if c.UserID == 0 {
		return fmt.Errorf("missing user id in claims: %w", jwt.ErrTokenInvalidClaims)
	}

	if c.OrgID == 0 {
		return fmt.Errorf("missing org id in claims: %w", jwt.ErrTokenInvalidClaims)
	}

	if c.OrgSubdomain == "" {
		return fmt.Errorf("missing org subdomain in claims: %w", jwt.ErrTokenInvalidClaims)
	}

	return nil
}

// GenerateJWT generates a new jwt token.
func GenerateJWT(appSecret string, userID, orgID int64, orgSubdomain string) (string, error) {
	claims := AppClaims{
		UserID:       userID,
		OrgID:        orgID,
		OrgSubdomain: orgSubdomain,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(SessionTTLDuration)),
		},
	}

	// create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// generate signed token using the secret signing key
	t, err := token.SignedString([]byte(appSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}

// ParseAndValidateJWT parses and validates the given jwt token string using the secret key.
func ParseAndValidateJWT(tokenString string, appSecret string) (*jwt.Token, *AppClaims, error) {
	claims := &AppClaims{}

	t, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(*jwt.Token) (any, error) {
			// return the secret signing key
			return []byte(appSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(), // make sure the exp claim is passed
	)
	if err != nil {
		return nil, nil, err
	}

	return t, claims, nil
}
