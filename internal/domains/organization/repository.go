package organization

import (
	"context"

	"github.com/camelhr/camelhr-api/internal/database"
)

//go:generate mockery --name=Repository --structname=RepositoryMock --inpackage --filename=repository_mock.go

type Repository interface {
	// GetOrganizationByID returns an organization by its ID.
	GetOrganizationByID(ctx context.Context, id int64) (Organization, error)

	// GetOrganizationBySubdomain returns an organization by its subdomain.
	GetOrganizationBySubdomain(ctx context.Context, subdomain string) (Organization, error)

	// GetOrganizationByName returns an organization by its name.
	GetOrganizationByName(ctx context.Context, name string) (Organization, error)

	// CreateOrganization creates a new organization.
	CreateOrganization(ctx context.Context, org Organization) (int64, error)

	// UpdateOrganization updates an organization.
	UpdateOrganization(ctx context.Context, org Organization) error

	// DeleteOrganization deletes an organization by its ID.
	DeleteOrganization(ctx context.Context, id int64) error

	// SuspendOrganization suspends an organization by its ID.
	SuspendOrganization(ctx context.Context, id int64, comment string) error

	// UnsuspendOrganization unsuspend an organization by its ID.
	UnsuspendOrganization(ctx context.Context, id int64, comment string) error

	// BlacklistOrganization blacklists an organization by its ID.
	BlacklistOrganization(ctx context.Context, id int64, comment string) error

	// UnblacklistOrganization unblacklist an organization by its ID.
	UnblacklistOrganization(ctx context.Context, id int64, comment string) error
}

type repository struct {
	db database.Database
}

func NewRepository(db database.Database) Repository {
	return &repository{db}
}

func (r *repository) GetOrganizationByID(ctx context.Context, id int64) (Organization, error) {
	var org Organization
	err := r.db.Get(ctx, &org, getOrganizationByIDQuery, id)

	return org, err
}

func (r *repository) GetOrganizationBySubdomain(ctx context.Context, subdomain string) (Organization, error) {
	var org Organization
	err := r.db.Get(ctx, &org, getOrganizationBySubdomainQuery, subdomain)

	return org, err
}

func (r *repository) GetOrganizationByName(ctx context.Context, name string) (Organization, error) {
	var org Organization
	err := r.db.Get(ctx, &org, getOrganizationByNameQuery, name)

	return org, err
}

func (r *repository) CreateOrganization(ctx context.Context, org Organization) (int64, error) {
	var id int64
	err := r.db.Exec(ctx, &id, createOrganizationQuery, org.Subdomain, org.Name)

	return id, err
}

func (r *repository) UpdateOrganization(ctx context.Context, org Organization) error {
	return r.db.Exec(ctx, nil, updateOrganizationQuery, org.ID, org.Subdomain, org.Name)
}

func (r *repository) DeleteOrganization(ctx context.Context, id int64) error {
	return r.db.Exec(ctx, nil, deleteOrganizationQuery, id)
}

func (r *repository) SuspendOrganization(ctx context.Context, id int64, comment string) error {
	return r.db.Exec(ctx, nil, suspendOrganizationQuery, id, comment)
}

func (r *repository) UnsuspendOrganization(ctx context.Context, id int64, comment string) error {
	return r.db.Exec(ctx, nil, unsuspendOrganizationQuery, id, comment)
}

func (r *repository) BlacklistOrganization(ctx context.Context, id int64, comment string) error {
	return r.db.Exec(ctx, nil, blacklistOrganizationQuery, id, comment)
}

func (r *repository) UnblacklistOrganization(ctx context.Context, id int64, comment string) error {
	return r.db.Exec(ctx, nil, unblacklistOrganizationQuery, id, comment)
}
