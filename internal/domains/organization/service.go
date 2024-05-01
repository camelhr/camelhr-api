package organization

import "context"

//go:generate mockery --name=Service --structname=ServiceMock --inpackage --filename=service_mock.go

type Service interface {
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

	// SuspendOrganization suspends an organization by its ID.
	SuspendOrganization(ctx context.Context, id int64, comment string) error

	// UnsuspendOrganization unsuspend an organization by its ID.
	UnsuspendOrganization(ctx context.Context, id int64, comment string) error

	// BlacklistOrganization blacklists an organization by its ID.
	BlacklistOrganization(ctx context.Context, id int64, comment string) error

	// UnblacklistOrganization unblacklist an organization by its ID.
	UnblacklistOrganization(ctx context.Context, id int64, comment string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetOrganizationByID(ctx context.Context, id int64) (Organization, error) {
	return s.repo.GetOrganizationByID(ctx, id)
}

func (s *service) GetOrganizationByName(ctx context.Context, name string) (Organization, error) {
	return s.repo.GetOrganizationByName(ctx, name)
}

func (s *service) CreateOrganization(ctx context.Context, org Organization) (int64, error) {
	return s.repo.CreateOrganization(ctx, org)
}

func (s *service) UpdateOrganization(ctx context.Context, org Organization) error {
	return s.repo.UpdateOrganization(ctx, org)
}

func (s *service) DeleteOrganization(ctx context.Context, id int64) error {
	return s.repo.DeleteOrganization(ctx, id)
}

func (s *service) SuspendOrganization(ctx context.Context, id int64, comment string) error {
	return s.repo.SuspendOrganization(ctx, id, comment)
}

func (s *service) UnsuspendOrganization(ctx context.Context, id int64, comment string) error {
	return s.repo.UnsuspendOrganization(ctx, id, comment)
}

func (s *service) BlacklistOrganization(ctx context.Context, id int64, comment string) error {
	return s.repo.BlacklistOrganization(ctx, id, comment)
}

func (s *service) UnblacklistOrganization(ctx context.Context, id int64, comment string) error {
	return s.repo.UnblacklistOrganization(ctx, id, comment)
}
