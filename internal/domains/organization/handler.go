package organization

import (
	"net/http"

	"github.com/camelhr/camelhr-api/internal/router/response"
)

type (
	organizationHandler struct {
	}

	OrganizationRequest struct {
		Name string `json:"name"`
	}

	OrganizationResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

func NewOrganizationHandler() *organizationHandler {
	return &organizationHandler{}
}

func (h *organizationHandler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	response.JSON(w, http.StatusOK, &OrganizationResponse{})
}

func (h *organizationHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *organizationHandler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *organizationHandler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
