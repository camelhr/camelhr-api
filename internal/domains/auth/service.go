package auth

import (
	"context"
	"errors"

	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/domains/session"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	// Register registers a new organization with owner.
	// The organization is set disabled by default.
	Register(ctx context.Context, email, password, subdomain, orgName string) error

	// Login logs in a user and returns a jwt token.
	Login(ctx context.Context, subdomain, email, password string) (string, error)

	// Logout logs out a user by deleting the session.
	Logout(ctx context.Context, userID, orgID int64) error
}

type service struct {
	appSecret      string
	transactor     database.Transactor
	orgService     organization.Service
	userService    user.Service
	sessionManager session.SessionManager
}

func NewService(
	appSecret string, transactor database.Transactor, orgService organization.Service,
	userService user.Service, sessionManager session.SessionManager,
) Service {
	return &service{
		appSecret:      appSecret,
		transactor:     transactor,
		orgService:     orgService,
		userService:    userService,
		sessionManager: sessionManager,
	}
}

var (
	ErrInvalidCredentials     = errors.New("email or password is invalid")
	ErrUserDisabled           = errors.New("user is disabled")
	ErrSubdomainAlreadyExists = errors.New("subdomain already exists")
)

func (s *service) Register(ctx context.Context, email, password, subdomain, orgName string) error {
	var err error

	// check if the subdomain already exists
	_, err = s.orgService.GetOrganizationBySubdomain(ctx, subdomain)
	if err == nil {
		return ErrSubdomainAlreadyExists
	} else if !base.IsNotFoundError(err) {
		return err
	}

	// create a new organization with owner. keep the organization in deleted state until verified
	err = s.transactor.WithTx(ctx, func(ctx context.Context) error {
		org, err := s.orgService.CreateOrganization(ctx, subdomain, orgName)
		if err != nil {
			return err
		}

		if _, err = s.userService.CreateOwner(ctx, org.ID, email, password); err != nil {
			return err
		}

		// delete the newly registered organization with a predefined comment
		// the organization & owner should be activated through backoffice upon verification
		return s.orgService.DeleteOrganization(ctx, org.ID, NewOrgDeleteComment)
	})

	return err
}

func (s *service) Login(ctx context.Context, subdomain, email, password string) (string, error) {
	org, err := s.orgService.GetOrganizationBySubdomain(ctx, subdomain)
	if err != nil {
		return "", err
	}

	u, err := s.userService.GetUserByOrgIDEmail(ctx, org.ID, email)
	if err != nil {
		if base.IsNotFoundError(err) {
			return "", ErrInvalidCredentials
		}

		return "", err
	}

	// prevent login for disabled user
	if u.DisabledAt != nil {
		return "", ErrUserDisabled
	}

	// compare the password with bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// generate a new jwt token with user and organization data
	jwtToken, err := GenerateJWT(s.appSecret, u.ID, org.ID, org.Subdomain)
	if err != nil {
		return "", err
	}

	// create session with the currently generated jwt token
	// if the session already exists, it will be updated with the new jwt token
	if err := s.sessionManager.CreateSession(
		ctx,
		u.ID,
		org.ID,
		jwtToken,
		ptrToString(u.APIToken),
		SessionTTLDuration,
	); err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (s *service) Logout(ctx context.Context, userID, orgID int64) error {
	return s.sessionManager.DeleteSession(ctx, userID, orgID)
}

func ptrToString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}
