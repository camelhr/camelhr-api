package organization

import (
	"context"

	"github.com/camelhr/camelhr-api/internal/database"
)

type Repository interface {
	// GetOrganizationByID returns an organization by its ID.
	GetOrganizationByID(ctx context.Context, id int64) (Organization, error)

	// GetOrganizationBySubdomain returns an organization by its subdomain.
	GetOrganizationBySubdomain(ctx context.Context, subdomain string) (Organization, error)

	// GetOrganizationByName returns an organization by its name.
	GetOrganizationByName(ctx context.Context, name string) (Organization, error)

	// CreateOrganization creates a new organization.
	CreateOrganization(ctx context.Context, subdomain string, name string) (Organization, error)

	// UpdateOrganization updates an organization.
	UpdateOrganization(ctx context.Context, id int64, name string) error

	// DeleteOrganization deletes an organization by its ID.
	DeleteOrganization(ctx context.Context, id int64) error

	// SuspendOrganization suspends an organization by its ID.
	SuspendOrganization(ctx context.Context, id int64, comment string) error

	// UnsuspendOrganization unsuspend an organization by its ID.
	UnsuspendOrganization(ctx context.Context, id int64, comment string) error

	// DisableOrganization disables an organization by its ID.
	DisableOrganization(ctx context.Context, id int64, comment string) error

	// EnableOrganization enable an organization by its ID.
	EnableOrganization(ctx context.Context, id int64, comment string) error
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

func (r *repository) CreateOrganization(ctx context.Context, subdomain string, name string) (Organization, error) {
	var org Organization
	err := r.db.Exec(ctx, &org, createOrganizationQuery, subdomain, name)

	return org, err
}

func (r *repository) UpdateOrganization(ctx context.Context, id int64, name string) error {
	return r.db.Exec(ctx, nil, updateOrganizationQuery, id, name)
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

func (r *repository) DisableOrganization(ctx context.Context, id int64, comment string) error {
	return r.db.Exec(ctx, nil, disableOrganizationQuery, id, comment)
}

func (r *repository) EnableOrganization(ctx context.Context, id int64, comment string) error {
	return r.db.Exec(ctx, nil, enableOrganizationQuery, id, comment)
}
