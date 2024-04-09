package organization

import (
	"net/http"

	"github.com/camelhr/camelhr-api/internal/router/request"
	"github.com/camelhr/camelhr-api/internal/router/response"
)

type (
	organizationHandler struct {
	}

	OrganizationRequest struct {
		Name string `json:"name"`
	}

	OrganizationResponse struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
)

func NewOrganizationHandler() *organizationHandler {
	return &organizationHandler{}
}

func (h *organizationHandler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	response.JSON(w, http.StatusOK, &OrganizationResponse{
		ID:   int64(1),
		Name: "Test Organization",
	})
}

func (h *organizationHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	var reqPayload OrganizationRequest
	if err := request.DecodeJSON(r.Body, &reqPayload); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}
	response.JSON(w, http.StatusCreated, &OrganizationResponse{
		ID:   int64(1),
		Name: reqPayload.Name,
	})
}

func (h *organizationHandler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *organizationHandler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
