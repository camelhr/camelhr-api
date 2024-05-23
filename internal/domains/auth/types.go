package auth

const (
	// APITokenBasicAuthPassword is the basic auth password when using api token.
	APITokenBasicAuthPassword = "api_token"
	// JWTCookieName is the name of the cookie that stores the jwt token.
	JWTCookieName = "jwt_session_id"
	// JWTMaxAgeSeconds is the max age of the jwt token in seconds.
	JWTMaxAgeSeconds = 2 * 60 * 60
)

type (
	// RegisterRequest represents the request payload for the register endpoint.
	RegisterRequest struct {
		Email     string `json:"email" validate:"email,required"`
		Password  string `json:"password" validate:"required,min=8,max=64"`
		Subdomain string `json:"organization_subdomain" validate:"required,alphanum,max=30"`
		OrgName   string `json:"organization_name" validate:"required,ascii,max=60"`
	}

	// LoginRequest represents the request payload for the login endpoint.
	LoginRequest struct {
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"required"`
	}
)
