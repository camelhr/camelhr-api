package organization

import "context"

//go:generate mockery --name=Service --structname=ServiceMock --inpackage --filename=service_mock.go

type Service interface {
	ListOrganizations(ctx context.Context) ([]Organization, error)
	GetOrganizationByID(ctx context.Context, id int64) (Organization, error)
	GetOrganizationByName(ctx context.Context, name string) (Organization, error)
	CreateOrganization(ctx context.Context, org Organization) (int64, error)
	UpdateOrganization(ctx context.Context, org Organization) error
	DeleteOrganization(ctx context.Context, id int64) error
}

type organizationService struct {
	repo Repository
}

func NewOrganizationService(repo Repository) Service {
	return &organizationService{
		repo: repo,
	}
}

func (s *organizationService) ListOrganizations(ctx context.Context) ([]Organization, error) {
	return s.repo.ListOrganizations(ctx)
}

func (s *organizationService) CreateOrganization(ctx context.Context, org Organization) (int64, error) {
	return s.repo.CreateOrganization(ctx, org)
}

func (s *organizationService) UpdateOrganization(ctx context.Context, org Organization) error {
	return s.repo.UpdateOrganization(ctx, org)
}

func (s *organizationService) DeleteOrganization(ctx context.Context, id int64) error {
	return s.repo.DeleteOrganization(ctx, id)
}

func (s *organizationService) GetOrganizationByID(ctx context.Context, id int64) (Organization, error) {
	return s.repo.GetOrganizationByID(ctx, id)
}

func (s *organizationService) GetOrganizationByName(ctx context.Context, name string) (Organization, error) {
	return s.repo.GetOrganizationByName(ctx, name)
}
