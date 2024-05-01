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

func (h *handler) GetOrganizationByID(w http.ResponseWriter, r *http.Request) {
	id, err := request.URLParamID(r, "orgID")
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	org, err := h.service.GetOrganizationByID(r.Context(), id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	resp := h.toResponse(org)
	response.JSON(w, http.StatusOK, resp)
}

func (h *handler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var reqPayload Request
	if err := request.DecodeJSON(r.Body, &reqPayload); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	org := h.toOrganization(reqPayload)

	id, err := h.service.CreateOrganization(r.Context(), org)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	org.ID = id
	resp := h.toResponse(org)

	response.JSON(w, http.StatusCreated, resp)
}

func (h *handler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	var reqPayload Request
	if err := request.DecodeJSON(r.Body, &reqPayload); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	org := h.toOrganization(reqPayload)

	if err := h.service.UpdateOrganization(r.Context(), org); err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	response.Empty(w, http.StatusOK)
}

func (h *handler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	id, err := request.URLParamID(r, "orgID")
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.DeleteOrganization(r.Context(), id); err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	response.Empty(w, http.StatusOK)
}

func (h *handler) toResponse(org Organization) *Response {
	return &Response{
		ID:          org.ID,
		Name:        org.Name,
		SuspendedAt: org.SuspendedAt,
		Timestamps:  org.Timestamps,
	}
}

func (h *handler) toOrganization(req Request) Organization {
	return Organization{
		Name: req.Name,
	}
}
