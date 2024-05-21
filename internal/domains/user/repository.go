package user

import (
	"context"

	"github.com/camelhr/camelhr-api/internal/database"
)

//go:generate mockery --name=Repository --structname=RepositoryMock --inpackage --filename=repository_mock.go

type Repository interface {
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
	CreateUser(ctx context.Context, orgID int64, email, passwordHash string, isOwner bool) (User, error)

	// ResetPassword resets the password of a user.
	ResetPassword(ctx context.Context, id int64, passwordHash string) error

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

type repository struct {
	db database.Database
}

func NewRepository(db database.Database) Repository {
	return &repository{db}
}

func (r *repository) GetUserByID(ctx context.Context, id int64) (User, error) {
	var user User
	err := r.db.Get(ctx, &user, getUserByIDQuery, id)

	return user, err
}

func (r *repository) GetUserByAPIToken(ctx context.Context, apiToken string) (User, error) {
	var user User
	err := r.db.Get(ctx, &user, getUserByAPITokenQuery, apiToken)

	return user, err
}

func (r *repository) GetUserByOrgSubdomainAPIToken(ctx context.Context, orgSubdomain, apiToken string) (User, error) {
	var user User
	err := r.db.Get(ctx, &user, getUserByOrgSubdomainAPITokenQuery, orgSubdomain, apiToken)

	return user, err
}

func (r *repository) GetUserByOrgIDEmail(ctx context.Context, orgID int64, email string) (User, error) {
	var user User
	err := r.db.Get(ctx, &user, getUserByOrgIDEmailQuery, orgID, email)

	return user, err
}

func (r *repository) GetUserByOrgSubdomainEmail(ctx context.Context, orgSubdomain, email string) (User, error) {
	var user User
	err := r.db.Get(ctx, &user, getUserByOrgSubdomainEmailQuery, orgSubdomain, email)

	return user, err
}

func (r *repository) CreateUser(ctx context.Context, orgID int64, email, passHash string, isOwner bool) (User, error) {
	var u User
	err := r.db.Exec(ctx, &u, createUserQuery, orgID, email, passHash, isOwner)

	return u, err
}

func (r *repository) ResetPassword(ctx context.Context, id int64, passwordHash string) error {
	return r.db.Exec(ctx, nil, resetPasswordQuery, id, passwordHash)
}

func (r *repository) DeleteUser(ctx context.Context, id int64) error {
	return r.db.Exec(ctx, nil, deleteUserQuery, id)
}

func (r *repository) DisableUser(ctx context.Context, id int64, comment string) error {
	return r.db.Exec(ctx, nil, disableUserQuery, id, comment)
}

func (r *repository) EnableUser(ctx context.Context, id int64, comment string) error {
	return r.db.Exec(ctx, nil, enableUserQuery, id, comment)
}

func (r *repository) GenerateAPIToken(ctx context.Context, id int64) error {
	return r.db.Exec(ctx, nil, generateAPITokenQuery, id)
}

func (r *repository) ResetAPIToken(ctx context.Context, id int64) error {
	return r.db.Exec(ctx, nil, resetAPITokenQuery, id)
}

func (r *repository) SetEmailVerified(ctx context.Context, id int64) error {
	return r.db.Exec(ctx, nil, setEmailVerifiedQuery, id)
}
