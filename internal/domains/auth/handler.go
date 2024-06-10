package auth

import (
	"errors"
	"net/http"

	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/web/request"
	"github.com/camelhr/camelhr-api/internal/web/response"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqPayload RegisterRequest
	if err := request.DecodeAndValidateJSON(r.Body, &reqPayload); err != nil {
		response.ErrorResponse(w, err)
		return
	}

	err := h.service.Register(ctx, reqPayload.Email, reqPayload.Password,
		reqPayload.Subdomain, reqPayload.OrgName)
	if err != nil {
		response.ErrorResponse(w, err)
		return
	}

	response.Empty(w, http.StatusCreated)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subdomain := request.URLParam(r, "subdomain")
	if err := organization.ValidateSubdomain(subdomain); err != nil {
		response.ErrorResponse(w, err)
		return
	}

	if err := r.ParseForm(); err != nil {
		response.ErrorResponse(w, base.WrapError(err, base.ErrorHTTPStatus(http.StatusBadRequest)))
		return
	}

	email := r.Form.Get("email")
	if err := user.ValidateEmail(email); err != nil {
		response.ErrorResponse(w, err)
		return
	}

	password := r.Form.Get("password")
	if err := user.ValidatePassword(password); err != nil {
		response.ErrorResponse(w, err)
		return
	}

	token, err := h.service.Login(ctx, subdomain, email, password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) || errors.Is(err, ErrUserDisabled) ||
			errors.Is(err, organization.ErrOrganizationDisabled) {
			response.ErrorResponse(w, base.WrapError(err, base.ErrorHTTPStatus(http.StatusUnauthorized)))
			return
		}

		response.ErrorResponse(w, err)

		return
	}

	response.SetCookie(w, JWTCookieName, token, JWTMaxAgeSeconds)
	response.Empty(w, http.StatusOK)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	response.RemoveCookie(w, JWTCookieName)
	response.Empty(w, http.StatusOK)
}
