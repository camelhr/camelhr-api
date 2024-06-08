package organization_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/camelhr/camelhr-api/internal/domains/auth"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
	"github.com/camelhr/camelhr-api/internal/web"
)

const (
	orgPathFormat = "/api/v1/subdomains/%s/organizations"
)

func (s *OrganizationTestSuite) TestHandlerIntegration_GetOrganizationBySubdomain() {
	s.Run("should return organization successfully", func() {
		s.T().Parallel()

		fakeOrg := fake.NewOrganization(s.DB)
		orgJSON, err := json.Marshal(toOrganizationResponse(fakeOrg.Organization))
		s.Require().NoError(err)

		// create a new request
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf(orgPathFormat, fakeOrg.Subdomain),
			nil,
		)
		s.Require().NoError(err)

		rr := httptest.NewRecorder()
		h := web.SetupRoutes(s.DB, s.Config)
		h.ServeHTTP(rr, req)

		// assert the response
		s.Require().Equal(http.StatusOK, rr.Code)
		s.JSONEq(string(orgJSON), rr.Body.String())

		// parse the response body
		var resp organization.Response
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		s.Require().NoError(err)

		// assert the response
		s.Equal(fakeOrg.Organization.ID, resp.ID)
		s.Equal(fakeOrg.Organization.Name, resp.Name)
		s.NotZero(resp.CreatedAt)
		s.NotZero(resp.UpdatedAt)
		s.Equal(resp.CreatedAt, resp.UpdatedAt)
	})
}

func (s *OrganizationTestSuite) TestHandlerIntegration_UpdateOrganization() {
	s.Run("should update organization successfully", func() {
		s.T().Parallel()

		fakeOrg := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, fakeOrg.Organization.ID)
		updatePayload := organization.Request{
			Name: randomOrganizationName(),
		}
		orgJSON, err := json.Marshal(updatePayload)
		s.Require().NoError(err)

		// create a new request
		req, err := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf(orgPathFormat, fakeOrg.Subdomain),
			strings.NewReader(string(orgJSON)),
		)
		s.Require().NoError(err)
		req.SetBasicAuth(*u.APIToken, auth.APITokenBasicAuthPassword)

		rr := httptest.NewRecorder()
		h := web.SetupRoutes(s.DB, s.Config)
		h.ServeHTTP(rr, req)

		// assert the response status code
		s.Require().Equal(http.StatusOK, rr.Code)
		s.Empty(rr.Body.String())

		// retrieve the updated organization
		result := fakeOrg.FetchLatest(s.DB)
		s.Require().NotNil(result)
		s.Equal(updatePayload.Name, result.Name)
	})
}

func (s *OrganizationTestSuite) TestHandlerIntegration_DeleteOrganization() {
	s.Run("should delete organization successfully", func() {
		s.T().Parallel()

		fakeOrg := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, fakeOrg.Organization.ID)

		// create a new request
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf(orgPathFormat, fakeOrg.Subdomain),
			nil,
		)
		s.Require().NoError(err)
		req.SetBasicAuth(*u.APIToken, auth.APITokenBasicAuthPassword)

		rr := httptest.NewRecorder()
		h := web.SetupRoutes(s.DB, s.Config)
		h.ServeHTTP(rr, req)

		// assert the response status code
		s.Require().Equal(http.StatusOK, rr.Code)
		s.Empty(rr.Body.String())

		// assert that the organization is deleted
		isDeleted := fakeOrg.IsDeleted(s.DB)
		s.True(isDeleted)
	})
}

func toOrganizationResponse(org organization.Organization) organization.Response {
	return organization.Response{
		ID:          org.ID,
		Subdomain:   org.Subdomain,
		Name:        org.Name,
		SuspendedAt: org.SuspendedAt,
		DisabledAt:  org.DisabledAt,
		Timestamps:  org.Timestamps,
	}
}
