package organization_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
	"github.com/camelhr/camelhr-api/internal/web"
)

func (s *OrganizationTestSuite) TestHandlerIntegration_GetOrganizationByID() {
	s.Run("should return organization successfully", func() {
		s.T().Parallel()

		fakeOrg := fake.NewOrganization(s.DB)
		orgJSON, err := json.Marshal(toOrganizationResponse(fakeOrg.Organization))
		s.Require().NoError(err)

		// create a new request
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/organizations/%d", fakeOrg.ID), nil)
		s.Require().NoError(err)

		rr := httptest.NewRecorder()
		h := web.SetupRoutes(s.DB)
		h.ServeHTTP(rr, req)

		// assert the response
		s.Equal(http.StatusOK, rr.Code)
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

func (s *OrganizationTestSuite) TestHandlerIntegration_CreateOrganization() {
	s.Run("should create organization successfully", func() {
		s.T().Parallel()

		orgReqPayload := organization.Request{
			Name: gofakeit.Name(),
		}
		orgJSON, err := json.Marshal(orgReqPayload)
		s.Require().NoError(err)

		// create a new request
		req, err := http.NewRequest(http.MethodPost, "/api/v1/organizations", strings.NewReader(string(orgJSON)))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		h := web.SetupRoutes(s.DB)
		h.ServeHTTP(rr, req)

		// assert the response status code
		s.Require().Equal(http.StatusCreated, rr.Code)

		// parse the response body
		var resp organization.Response
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		s.Require().NoError(err)

		// assert the response
		s.NotEmpty(resp.Name)
		s.Equal(orgReqPayload.Name, resp.Name)
		s.NotZero(resp.ID)
	})
}

func (s *OrganizationTestSuite) TestHandlerIntegration_UpdateOrganization() {
	s.Run("should update organization successfully", func() {
		s.T().Parallel()

		fakeOrg := fake.NewOrganization(s.DB)
		orgReqPayload := organization.Request{
			Name: gofakeit.Name(),
		}
		orgJSON, err := json.Marshal(orgReqPayload)
		s.Require().NoError(err)

		// create a new request
		req, err := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("/api/v1/organizations/%d", fakeOrg.ID),
			strings.NewReader(string(orgJSON)),
		)
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		h := web.SetupRoutes(s.DB)
		h.ServeHTTP(rr, req)

		// assert the response status code
		s.Require().Equal(http.StatusOK, rr.Code)
	})
}

func (s *OrganizationTestSuite) TestHandlerIntegration_DeleteOrganization() {
	s.Run("should delete organization successfully", func() {
		s.T().Parallel()

		fakeOrg := fake.NewOrganization(s.DB)

		// create a new request
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/organizations/%d", fakeOrg.ID), nil)
		s.Require().NoError(err)

		rr := httptest.NewRecorder()
		h := web.SetupRoutes(s.DB)
		h.ServeHTTP(rr, req)

		// assert the response status code
		s.Require().Equal(http.StatusOK, rr.Code)
	})
}

func toOrganizationResponse(org organization.Organization) organization.Response {
	return organization.Response{
		ID:            org.ID,
		Name:          org.Name,
		SuspendedAt:   org.SuspendedAt,
		BlacklistedAt: org.BlacklistedAt,
		Timestamps:    org.Timestamps,
	}
}