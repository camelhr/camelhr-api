package organization

import (
	"net/http"

	"github.com/camelhr/camelhr-api/internal/web/request"
	"github.com/camelhr/camelhr-api/internal/web/response"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service}
}

func (h *handler) GetOrganizationBySubdomain(w http.ResponseWriter, r *http.Request) {
	subdomain := request.URLParam(r, "subdomain")
	if err := ValidateSubdomain(subdomain); err != nil {
		response.ErrorResponse(w, err)
		return
	}

	org, err := h.service.GetOrganizationBySubdomain(r.Context(), subdomain)
	if err != nil {
		response.ErrorResponse(w, err)
		return
	}

	resp := h.toResponse(org)
	response.JSON(w, http.StatusOK, resp)
}

func (h *handler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	subdomain := request.URLParam(r, "subdomain")
	if err := ValidateSubdomain(subdomain); err != nil {
		response.ErrorResponse(w, err)
		return
	}

	var reqPayload UpdateRequest
	if err := request.DecodeAndValidateJSON(r.Body, &reqPayload); err != nil {
		response.ErrorResponse(w, err)

		return
	}

	org, err := h.service.GetOrganizationBySubdomain(r.Context(), subdomain)
	if err != nil {
		response.ErrorResponse(w, err)
		return
	}

	if err := h.service.UpdateOrganization(r.Context(), org.ID, reqPayload.Name); err != nil {
		response.ErrorResponse(w, err)
		return
	}

	response.Empty(w, http.StatusOK)
}

func (h *handler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	subdomain := request.URLParam(r, "subdomain")
	if err := ValidateSubdomain(subdomain); err != nil {
		response.ErrorResponse(w, err)
		return
	}

	org, err := h.service.GetOrganizationBySubdomain(r.Context(), subdomain)
	if err != nil {
		response.ErrorResponse(w, err)
		return
	}

	var reqPayload DeleteRequest
	if err := request.DecodeAndValidateJSON(r.Body, &reqPayload); err != nil {
		response.ErrorResponse(w, err)

		return
	}

	if err := ValidateComment(reqPayload.Comment); err != nil {
		response.ErrorResponse(w, err)
		return
	}

	if err := h.service.DeleteOrganization(r.Context(), org.ID, reqPayload.Comment); err != nil {
		response.ErrorResponse(w, err)
		return
	}

	response.Empty(w, http.StatusOK)
}

func (h *handler) toResponse(org Organization) *Response {
	return &Response{
		ID:          org.ID,
		Subdomain:   org.Subdomain,
		Name:        org.Name,
		SuspendedAt: org.SuspendedAt,
		Timestamps:  org.Timestamps,
	}
}
