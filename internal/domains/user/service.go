package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/domains/session"
	"golang.org/x/crypto/bcrypt"
)

// Service is a service for managing users.
type Service interface {
	// GetUserByID returns a user by its ID.
	GetUserByID(ctx context.Context, id int64) (User, error)

	// GetUserByAPIToken returns a user by its API token.
	GetUserByAPIToken(ctx context.Context, apiToken string) (User, error)

	// GetUserByOrgSubdomainAPIToken returns a user of organization by its org subdomain and api token.
	GetUserByOrgSubdomainAPIToken(ctx context.Context, orgSubdomain, apiToken string) (User, error)

	// GetUserByOrgIDEmail returns a user of organization by its org id and email.
	GetUserByOrgIDEmail(ctx context.Context, orgID int64, email string) (User, error)

	// GetUserByOrgSubdomainEmail returns a user of organization by its org subdomain and email.
	GetUserByOrgSubdomainEmail(ctx context.Context, orgSubdomain, email string) (User, error)

	// CreateUser creates a new user.
	CreateUser(ctx context.Context, orgID int64, email, password string) (User, error)

	// CreateOwner creates a new owner user.
	CreateOwner(ctx context.Context, orgID int64, email, password string) (User, error)

	// ResetPassword resets the password of a user.
	ResetPassword(ctx context.Context, id int64, newPassword string) error

	// DeleteUser deletes a user by its ID.
	// This also deletes the user session.
	DeleteUser(ctx context.Context, id int64, comment string) error

	// DisableUser disables a user by its ID.
	// This also deletes the user session.
	DisableUser(ctx context.Context, id int64, comment string) error

	// EnableUser enables a user by its ID.
	EnableUser(ctx context.Context, id int64, comment string) error

	// GenerateAPIToken generates a new API token for a user.
	GenerateAPIToken(ctx context.Context, id int64) error

	// ResetAPIToken resets the API token of a user.
	ResetAPIToken(ctx context.Context, id int64) error

	// SetEmailVerified sets the email_verified flag of a user.
	SetEmailVerified(ctx context.Context, id int64) error
}

var ErrUserIsOwner = errors.New("operation not allowed. user is owner")

type service struct {
	repo           Repository
	sessionManager session.SessionManager
}

// NewService creates a new user service.
func NewService(repo Repository, sessionManager session.SessionManager) *service {
	return &service{repo, sessionManager}
}

func (s *service) GetUserByID(ctx context.Context, id int64) (User, error) {
	u, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, base.NewNotFoundError("user not found for the given id")
		}

		return User{}, err
	}

	return u, nil
}

func (s *service) GetUserByAPIToken(ctx context.Context, apiToken string) (User, error) {
	u, err := s.repo.GetUserByAPIToken(ctx, apiToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, base.NewNotFoundError("user not found for the given api-token")
		}

		return User{}, err
	}

	return u, nil
}

func (s *service) GetUserByOrgSubdomainAPIToken(ctx context.Context, orgSubdomain, apiToken string) (User, error) {
	if err := organization.ValidateSubdomain(orgSubdomain); err != nil {
		return User{}, err
	}

	u, err := s.repo.GetUserByOrgSubdomainAPIToken(ctx, orgSubdomain, apiToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, base.NewNotFoundError("user not found for the given org-subdomain and api-token")
		}

		return User{}, err
	}

	return u, nil
}

func (s *service) GetUserByOrgIDEmail(ctx context.Context, orgID int64, email string) (User, error) {
	if err := ValidateEmail(email); err != nil {
		return User{}, err
	}

	u, err := s.repo.GetUserByOrgIDEmail(ctx, orgID, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, base.NewNotFoundError("user not found for the given org-id and email")
		}

		return User{}, err
	}

	return u, nil
}

func (s *service) GetUserByOrgSubdomainEmail(ctx context.Context, orgSubdomain, email string) (User, error) {
	if err := organization.ValidateSubdomain(orgSubdomain); err != nil {
		return User{}, err
	}

	if err := ValidateEmail(email); err != nil {
		return User{}, err
	}

	u, err := s.repo.GetUserByOrgSubdomainEmail(ctx, orgSubdomain, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, base.NewNotFoundError("user not found for the given org-subdomain and email")
		}

		return User{}, err
	}

	return u, nil
}

func (s *service) CreateUser(ctx context.Context, orgID int64, email, password string) (User, error) {
	if err := ValidateEmail(email); err != nil {
		return User{}, err
	}

	if err := ValidatePassword(password); err != nil {
		return User{}, err
	}

	passwordHash, err := s.bcryptPassword(password)
	if err != nil {
		return User{}, err
	}

	return s.repo.CreateUser(ctx, orgID, email, passwordHash, false)
}

func (s *service) CreateOwner(ctx context.Context, orgID int64, email, password string) (User, error) {
	if err := ValidateEmail(email); err != nil {
		return User{}, err
	}

	if err := ValidatePassword(password); err != nil {
		return User{}, err
	}

	passwordHash, err := s.bcryptPassword(password)
	if err != nil {
		return User{}, err
	}

	return s.repo.CreateUser(ctx, orgID, email, passwordHash, true)
}

func (s *service) ResetPassword(ctx context.Context, id int64, newPassword string) error {
	if err := ValidatePassword(newPassword); err != nil {
		return err
	}

	passwordHash, err := s.bcryptPassword(newPassword)
	if err != nil {
		return err
	}

	return s.repo.ResetPassword(ctx, id, passwordHash)
}

func (s *service) DeleteUser(ctx context.Context, id int64, comment string) error {
	if err := ValidateComment(comment); err != nil {
		return err
	}

	u, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return base.NewNotFoundError("user not found for the given id")
		}

		return err
	}

	if u.IsOwner {
		return ErrUserIsOwner
	}

	if err := s.repo.DeleteUser(ctx, id, comment); err != nil {
		return err
	}

	return s.sessionManager.DeleteSession(ctx, u.ID, u.OrganizationID)
}

func (s *service) DisableUser(ctx context.Context, id int64, comment string) error {
	if err := ValidateComment(comment); err != nil {
		return err
	}

	u, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return base.NewNotFoundError("user not found for the given id")
		}

		return err
	}

	if u.IsOwner {
		return ErrUserIsOwner
	}

	if err := s.repo.DisableUser(ctx, id, comment); err != nil {
		return err
	}

	return s.sessionManager.DeleteSession(ctx, u.ID, u.OrganizationID)
}

func (s *service) EnableUser(ctx context.Context, id int64, comment string) error {
	if err := ValidateComment(comment); err != nil {
		return err
	}

	return s.repo.EnableUser(ctx, id, comment)
}

func (s *service) GenerateAPIToken(ctx context.Context, id int64) error {
	return s.repo.GenerateAPIToken(ctx, id)
}

func (s *service) ResetAPIToken(ctx context.Context, id int64) error {
	return s.repo.ResetAPIToken(ctx, id)
}

func (s *service) SetEmailVerified(ctx context.Context, id int64) error {
	return s.repo.SetEmailVerified(ctx, id)
}

// bcryptPassword hashes a password using bcrypt.
func (s *service) bcryptPassword(password string) (string, error) {
	passwordBytes := []byte(password)

	// generate a bcrypt hash of the password. this automatically generates a new salt for each password
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// verify that the password and the hashed password match
	// this is not typically necessary immediately after hashing the password
	err = bcrypt.CompareHashAndPassword(passwordHashBytes, passwordBytes)
	if err != nil {
		return "", err
	}

	// return the hashed password as a string
	return string(passwordHashBytes), nil
}
