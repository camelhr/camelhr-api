package organization

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/camelhr/camelhr-api/internal/base"
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
	if err := validateSubdomain(subdomain); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	org, err := h.service.GetOrganizationBySubdomain(r.Context(), subdomain)
	if err != nil {
		var notFoundErr *base.NotFoundError
		if errors.As(err, &notFoundErr) {
			response.ErrorResponse(w, http.StatusNotFound, notFoundErr)
			return
		}

		response.ErrorResponse(w, http.StatusInternalServerError, err)

		return
	}

	resp := h.toResponse(org)
	response.JSON(w, http.StatusOK, resp)
}

func (h *handler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	subdomain := request.URLParam(r, "subdomain")
	if err := validateSubdomain(subdomain); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	var reqPayload Request
	if err := request.DecodeAndValidateJSON(r.Body, &reqPayload); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	org, err := h.service.GetOrganizationBySubdomain(r.Context(), subdomain)
	if err != nil {
		var notFoundErr *base.NotFoundError
		if errors.As(err, &notFoundErr) {
			response.ErrorResponse(w, http.StatusNotFound, notFoundErr)
			return
		}

		response.ErrorResponse(w, http.StatusInternalServerError, err)

		return
	}

	if err := h.service.UpdateOrganization(r.Context(), org.ID, reqPayload.Name); err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	response.Empty(w, http.StatusOK)
}

func (h *handler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	subdomain := request.URLParam(r, "subdomain")
	if err := validateSubdomain(subdomain); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	org, err := h.service.GetOrganizationBySubdomain(r.Context(), subdomain)
	if err != nil {
		var notFoundErr *base.NotFoundError
		if errors.As(err, &notFoundErr) {
			response.ErrorResponse(w, http.StatusNotFound, notFoundErr)
			return
		}

		response.ErrorResponse(w, http.StatusInternalServerError, err)

		return
	}

	if err := h.service.DeleteOrganization(r.Context(), org.ID); err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	response.Empty(w, http.StatusOK)
}

func (h *handler) toResponse(org Organization) *Response {
	return &Response{
		ID:            org.ID,
		Subdomain:     org.Subdomain,
		Name:          org.Name,
		SuspendedAt:   org.SuspendedAt,
		BlacklistedAt: org.BlacklistedAt,
		Timestamps:    org.Timestamps,
	}
}

func validateSubdomain(subdomain string) error {
	const allowedMaxLength = 30

	// validate that subdomain is not empty
	if subdomain == "" {
		return base.NewAPIError("subdomain should not be empty")
	}

	// validate that subdomain length is less than or equal to 30
	if len(subdomain) > allowedMaxLength {
		return base.NewAPIError("subdomain length must be less than or equal to 30")
	}

	// validate that subdomain is alphanumeric
	match, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, subdomain)
	if err != nil || !match {
		return base.NewAPIError("subdomain must be alphanumeric")
	}

	return nil
}
