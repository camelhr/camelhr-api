package user

import (
	"context"
	"errors"

	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"golang.org/x/crypto/bcrypt"
)

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
	DeleteUser(ctx context.Context, id int64) error

	// DisableUser disables a user by its ID.
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

var ErrUserCommentMissing = errors.New("comment is missing")

type service struct {
	repo Repository
}

// NewService creates a new user service.
func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) GetUserByID(ctx context.Context, id int64) (User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *service) GetUserByAPIToken(ctx context.Context, apiToken string) (User, error) {
	return s.repo.GetUserByAPIToken(ctx, apiToken)
}

func (s *service) GetUserByOrgSubdomainAPIToken(ctx context.Context, orgSubdomain, apiToken string) (User, error) {
	if err := organization.ValidateSubdomain(orgSubdomain); err != nil {
		return User{}, err
	}

	return s.repo.GetUserByOrgSubdomainAPIToken(ctx, orgSubdomain, apiToken)
}

func (s *service) GetUserByOrgIDEmail(ctx context.Context, orgID int64, email string) (User, error) {
	if err := ValidateEmail(email); err != nil {
		return User{}, err
	}

	return s.repo.GetUserByOrgIDEmail(ctx, orgID, email)
}

func (s *service) GetUserByOrgSubdomainEmail(ctx context.Context, orgSubdomain, email string) (User, error) {
	if err := organization.ValidateSubdomain(orgSubdomain); err != nil {
		return User{}, err
	}

	if err := ValidateEmail(email); err != nil {
		return User{}, err
	}

	return s.repo.GetUserByOrgSubdomainEmail(ctx, orgSubdomain, email)
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

func (s *service) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *service) DisableUser(ctx context.Context, id int64, comment string) error {
	if comment == "" {
		return ErrUserCommentMissing
	}

	return s.repo.DisableUser(ctx, id, comment)
}

func (s *service) EnableUser(ctx context.Context, id int64, comment string) error {
	if comment == "" {
		return ErrUserCommentMissing
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
