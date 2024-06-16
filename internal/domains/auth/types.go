package auth

import "time"

const (
	// APITokenBasicAuthPassword is the basic auth password when using api token.
	APITokenBasicAuthPassword = "api_token"
	// JWTCookieName is the name of the cookie that stores the jwt token.
	JWTCookieName = "jwt_session_id"

	// DefaultSessionTTL is the time duration to keep the session alive.
	DefaultSessionTTL = 24 * time.Hour

	// RememberMeSessionTTL is the time duration to keep the session alive when remember me is enabled.
	RememberMeSessionTTL = 8 * DefaultSessionTTL

	// NewOrgDeleteComment is the comment message to identify newly registered organizations
	// that are pending verification.
	NewOrgDeleteComment = "deletion_reason: new_unverified_organization"
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
