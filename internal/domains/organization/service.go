package organization

import (
	"context"
	"database/sql"
	"errors"

	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/camelhr-api/internal/domains/session"
)

type Service interface {
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

type service struct {
	repo           Repository
	sessionManager session.SessionManager
}

func NewService(repo Repository, sessionManager session.SessionManager) Service {
	return &service{repo, sessionManager}
}

func (s *service) GetOrganizationByID(ctx context.Context, id int64) (Organization, error) {
	o, err := s.repo.GetOrganizationByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Organization{}, base.NewNotFoundError("organization not found for the given id")
	}

	return o, err
}

func (s *service) GetOrganizationBySubdomain(ctx context.Context, subdomain string) (Organization, error) {
	if err := ValidateSubdomain(subdomain); err != nil {
		return Organization{}, err
	}

	o, err := s.repo.GetOrganizationBySubdomain(ctx, subdomain)
	if errors.Is(err, sql.ErrNoRows) {
		return Organization{}, base.NewNotFoundError("organization not found for the given subdomain")
	}

	return o, err
}

func (s *service) GetOrganizationByName(ctx context.Context, name string) (Organization, error) {
	if err := ValidateOrgName(name); err != nil {
		return Organization{}, err
	}

	o, err := s.repo.GetOrganizationByName(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		return Organization{}, base.NewNotFoundError("organization not found for the given name")
	}

	return o, err
}

func (s *service) CreateOrganization(ctx context.Context, subdomain string, name string) (Organization, error) {
	if err := ValidateSubdomain(subdomain); err != nil {
		return Organization{}, err
	}

	if err := ValidateOrgName(name); err != nil {
		return Organization{}, err
	}

	return s.repo.CreateOrganization(ctx, subdomain, name)
}

func (s *service) UpdateOrganization(ctx context.Context, id int64, name string) error {
	if err := ValidateOrgName(name); err != nil {
		return err
	}

	return s.repo.UpdateOrganization(ctx, id, name)
}

func (s *service) DeleteOrganization(ctx context.Context, id int64) error {
	if err := s.repo.DeleteOrganization(ctx, id); err != nil {
		return err
	}

	return s.sessionManager.DeleteAllOrgSessions(ctx, id)
}

func (s *service) SuspendOrganization(ctx context.Context, id int64, comment string) error {
	return s.repo.SuspendOrganization(ctx, id, comment)
}

func (s *service) UnsuspendOrganization(ctx context.Context, id int64, comment string) error {
	return s.repo.UnsuspendOrganization(ctx, id, comment)
}

func (s *service) DisableOrganization(ctx context.Context, id int64, comment string) error {
	return s.repo.DisableOrganization(ctx, id, comment)
}

func (s *service) EnableOrganization(ctx context.Context, id int64, comment string) error {
	return s.repo.EnableOrganization(ctx, id, comment)
}
