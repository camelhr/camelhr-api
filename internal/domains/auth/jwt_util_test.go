package auth_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWT(t *testing.T) {
	t.Parallel()

	t.Run("should generate a valid jwt token", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		appSecret := gofakeit.UUID()

		// create random user id, org id and org subdomain
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		orgSubdomain := gofakeit.Username()

		// generate jwt token
		token, err := auth.GenerateJWT(appSecret, userID, orgID, orgSubdomain)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// verify the claims in the token
		parsedToken, err := jwt.ParseWithClaims(token, &auth.AppClaims{}, func(*jwt.Token) (any, error) {
			return []byte(appSecret), nil
		})
		require.NoError(t, err)
		require.True(t, parsedToken.Valid)
		require.NotNil(t, parsedToken.Claims)
		require.IsType(t, &auth.AppClaims{}, parsedToken.Claims)

		appClaims, _ := parsedToken.Claims.(*auth.AppClaims)
		assert.Equal(t, userID, appClaims.UserID)
		assert.Equal(t, orgID, appClaims.OrgID)
		assert.Equal(t, orgSubdomain, appClaims.OrgSubdomain)

		now := time.Now()
		exp, err := parsedToken.Claims.GetExpirationTime()
		require.NoError(t, err)
		require.NotNil(t, exp)
		assert.WithinDuration(t, now, exp.Time, 24*time.Hour) //nolint:testifylint // false positive

		iat, err := parsedToken.Claims.GetIssuedAt()
		require.NoError(t, err)
		require.NotNil(t, iat)
		assert.WithinDuration(t, now, iat.Time, 1*time.Minute)
	})
}

func TestParseAndValidateJWT(t *testing.T) { //nolint:maintidx // test function
	t.Parallel()

	t.Run("should parse and validate a valid jwt token", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		appSecret := gofakeit.UUID()

		// create random user id, org id and org subdomain
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		orgSubdomain := gofakeit.Username()

		// generate jwt token
		const hoursValid = 24
		claims := auth.AppClaims{
			UserID:       userID,
			OrgID:        orgID,
			OrgSubdomain: orgSubdomain,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(hoursValid * time.Hour)),
			},
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(appSecret))
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// parse and validate the token
		parsedToken, parsedClaims, err := auth.ParseAndValidateJWT(token, appSecret)
		require.NoError(t, err)
		require.True(t, parsedToken.Valid)
		require.Equal(t, parsedToken.Claims, parsedClaims) // claims from token and parsed claims should be same

		// verify the claims
		require.NotNil(t, parsedClaims)
		assert.Equal(t, userID, parsedClaims.UserID)
		assert.Equal(t, orgID, parsedClaims.OrgID)
		assert.Equal(t, orgSubdomain, parsedClaims.OrgSubdomain)
	})

	t.Run("should return error if jwt is expired", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		appSecret := gofakeit.UUID()

		// create random user id, org id and org subdomain
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		orgSubdomain := gofakeit.Username()

		// generate expired jwt token
		const expiry = -auth.JWTMaxAgeSeconds + 120
		claims := auth.AppClaims{
			UserID:       userID,
			OrgID:        orgID,
			OrgSubdomain: orgSubdomain,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry * time.Second)),
			},
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(appSecret))
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// parse and validate the token
		parsedToken, parsedClaims, err := auth.ParseAndValidateJWT(token, appSecret)
		require.Error(t, err)
		require.ErrorIs(t, err, jwt.ErrTokenExpired)
		require.ErrorContains(t, err, "token is expired")
		assert.Nil(t, parsedToken)
		assert.Nil(t, parsedClaims)
	})

	t.Run("should return error when app claims are not present in token", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		appSecret := gofakeit.UUID()

		// generate token without app claims, include only standard claims
		const hoursValid = 24
		claims := jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(hoursValid * time.Hour)),
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(appSecret))
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// parse and validate the token
		parsedToken, parsedClaims, err := auth.ParseAndValidateJWT(token, appSecret)
		require.Error(t, err)
		require.ErrorIs(t, err, jwt.ErrTokenInvalidClaims)
		require.ErrorContains(t, err, "missing user id in claims")
		assert.Nil(t, parsedToken)
		assert.Nil(t, parsedClaims)
	})

	t.Run("should return error when app claims fields are missing", func(t *testing.T) {
		t.Parallel()

		// table test
		tests := []struct {
			testName  string
			userID    int64
			orgID     int64
			orgSub    string
			errString string
		}{
			{
				testName:  "missing user id",
				userID:    0,
				orgID:     gofakeit.Int64(),
				orgSub:    gofakeit.Word(),
				errString: "missing user id in claims",
			},
			{
				testName:  "missing org id",
				userID:    gofakeit.Int64(),
				orgID:     0,
				orgSub:    gofakeit.Word(),
				errString: "missing org id in claims",
			},
			{
				testName:  "missing org subdomain",
				userID:    gofakeit.Int64(),
				orgID:     gofakeit.Int64(),
				orgSub:    "",
				errString: "missing org subdomain in claims",
			},
		}

		for _, tt := range tests {
			// avoid loop closure issue by defining the variables here
			userID := tt.userID
			orgID := tt.orgID
			orgSub := tt.orgSub
			errString := tt.errString

			t.Run(tt.testName, func(t *testing.T) {
				t.Parallel()

				// create a random secret key
				appSecret := gofakeit.UUID()

				// generate jwt token
				const hoursValid = 24
				claims := auth.AppClaims{
					UserID:       userID,
					OrgID:        orgID,
					OrgSubdomain: orgSub,
					RegisteredClaims: jwt.RegisteredClaims{
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(hoursValid * time.Hour)),
					},
				}
				token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(appSecret))
				require.NoError(t, err)
				require.NotEmpty(t, token)

				// parse and validate the token
				parsedToken, parsedClaims, err := auth.ParseAndValidateJWT(token, appSecret)
				require.Error(t, err)
				require.ErrorIs(t, err, jwt.ErrTokenInvalidClaims)
				require.ErrorContains(t, err, errString)
				assert.Nil(t, parsedToken)
				assert.Nil(t, parsedClaims)
			})
		}
	})

	t.Run("should return error when exp claim is missing", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		appSecret := gofakeit.UUID()

		// generate token without exp claims
		claims := auth.AppClaims{
			UserID:       gofakeit.Int64(),
			OrgID:        gofakeit.Int64(),
			OrgSubdomain: gofakeit.Word(),
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt: jwt.NewNumericDate(time.Now()),
			},
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(appSecret))
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// parse and validate the token
		parsedToken, parsedClaims, err := auth.ParseAndValidateJWT(token, appSecret)
		require.Error(t, err)
		require.ErrorIs(t, err, jwt.ErrTokenInvalidClaims)
		require.ErrorContains(t, err, "exp claim is required")
		assert.Nil(t, parsedToken)
		assert.Nil(t, parsedClaims)
	})

	t.Run("should return error when method(alg) is not HS256", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		appSecret := gofakeit.UUID()

		// generate token without eat claims
		const hoursValid = 24
		claims := auth.AppClaims{
			UserID:       gofakeit.Int64(),
			OrgID:        gofakeit.Int64(),
			OrgSubdomain: gofakeit.Word(),
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(hoursValid * time.Hour)),
			},
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS384, claims).SignedString([]byte(appSecret))
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// parse and validate the token
		parsedToken, parsedClaims, err := auth.ParseAndValidateJWT(token, appSecret)
		require.Error(t, err)
		require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
		require.ErrorContains(t, err, "signing method HS384 is invalid")
		assert.Nil(t, parsedToken)
		assert.Nil(t, parsedClaims)
	})

	t.Run("should return error when method(alg) is none", func(t *testing.T) {
		t.Parallel()

		// generate token without eat claims
		const hoursValid = 24
		claims := auth.AppClaims{
			UserID:       gofakeit.Int64(),
			OrgID:        gofakeit.Int64(),
			OrgSubdomain: gofakeit.Word(),
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(hoursValid * time.Hour)),
			},
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// parse and validate the token
		parsedToken, parsedClaims, err := auth.ParseAndValidateJWT(token, "")
		require.Error(t, err)
		require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
		require.ErrorContains(t, err, "signing method none is invalid")
		assert.Nil(t, parsedToken)
		assert.Nil(t, parsedClaims)
	})

	t.Run("should return an error when the token is invalid", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		appSecret := gofakeit.UUID()

		// create random user id, org id and org subdomain
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		orgSubdomain := gofakeit.Username()

		// generate jwt token
		token, err := auth.GenerateJWT(appSecret, userID, orgID, orgSubdomain)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// modify the token
		token += "invalid"

		// parse and validate the token
		_, _, err = auth.ParseAndValidateJWT(token, appSecret)
		require.Error(t, err)
	})
}
