package auth

const (
	// APITokenPassword is the basic auth password when using api token.
	APITokenPassword = "api_token"
	// JWTCookieName is the name of the cookie that stores the jwt token.
	JWTCookieName = "jwt_session_id"
	// JWTMaxAgeSeconds is the max age of the jwt token in seconds.
	JWTMaxAgeSeconds = 2 * 60 * 60
)
