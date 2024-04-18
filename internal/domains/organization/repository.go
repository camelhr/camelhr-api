package organization

import (
	"context"

	"github.com/camelhr/camelhr-api/internal/database"
)

//go:generate mockery --name=Repository --structname=RepositoryMock --inpackage --filename=repository_mock.go

type Repository interface {
	// GetOrganizationByID returns an organization by its ID.
	GetOrganizationByID(ctx context.Context, id int64) (Organization, error)

	// GetOrganizationByName returns an organization by its name.
	GetOrganizationByName(ctx context.Context, name string) (Organization, error)

	// CreateOrganization creates a new organization.
	CreateOrganization(ctx context.Context, org Organization) (int64, error)

	// UpdateOrganization updates an organization.
	UpdateOrganization(ctx context.Context, org Organization) error

	// DeleteOrganization deletes an organization by its ID.
	DeleteOrganization(ctx context.Context, id int64) error
}

type organizationRepository struct {
	db database.Database
}

func NewOrganizationRepository(db database.Database) Repository {
	return &organizationRepository{db}
}

func (r *organizationRepository) GetOrganizationByID(ctx context.Context, id int64) (Organization, error) {
	var org Organization
	err := r.db.Get(ctx, &org, getOrganizationByIDQuery, id)

	return org, err
}

func (r *organizationRepository) GetOrganizationByName(ctx context.Context, name string) (Organization, error) {
	var org Organization
	err := r.db.Get(ctx, &org, getOrganizationByNameQuery, name)

	return org, err
}

func (r *organizationRepository) CreateOrganization(ctx context.Context, org Organization) (int64, error) {
	var id int64
	err := r.db.Exec(ctx, &id, createOrganizationQuery, org.Name)

	return id, err
}

func (r *organizationRepository) UpdateOrganization(ctx context.Context, org Organization) error {
	return r.db.Exec(ctx, nil, updateOrganizationQuery, org.Name)
}

func (r *organizationRepository) DeleteOrganization(ctx context.Context, id int64) error {
	return r.db.Exec(ctx, nil, deleteOrganizationQuery, id)
}
