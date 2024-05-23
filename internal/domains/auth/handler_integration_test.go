package auth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/auth"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
	"github.com/camelhr/camelhr-api/internal/web"
)

const (
	loginPathFormat = "/api/v1/subdomains/%s/auth/login"
)

func (s *AuthTestSuite) TestHandlerIntegration_Register() {
	s.Run("should register successfully", func() {
		s.T().Parallel()

		req, err := http.NewRequest(
			http.MethodPost,
			registerPath,
			strings.NewReader(
				fmt.Sprintf(`{"email":"%s","password":"%s","organization_name":"%s","organization_subdomain":"%s"}`,
					gofakeit.Email(), validPassword, gofakeit.LetterN(50), gofakeit.LetterN(20)),
			),
		)
		s.Require().NoError(err)

		rr := httptest.NewRecorder()
		h := web.SetupRoutes(s.DB, s.Config)
		h.ServeHTTP(rr, req)

		// assert the response
		s.Require().Equal(http.StatusCreated, rr.Code)
		s.Empty(rr.Body.String())
	})
}

func (s *AuthTestSuite) TestHandlerIntegration_Login() {
	s.Run("should login successfully", func() {
		s.T().Parallel()

		// create an organization and a user
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserPassword(validPassword))

		// create url-encoded form data
		form := url.Values{}
		form.Add("email", u.Email)
		form.Add("password", validPassword)
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(loginPathFormat, o.Subdomain),
			strings.NewReader(form.Encode()))
		s.Require().NoError(err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		h := web.SetupRoutes(s.DB, s.Config)
		h.ServeHTTP(rr, req)

		// assert the response
		s.Require().Equal(http.StatusOK, rr.Code)
		s.Empty(rr.Body.String())
		s.Contains(rr.Header().Get("Set-Cookie"), auth.JWTCookieName)
	})
}

func (s *AuthTestSuite) TestHandlerIntegration_Logout() {
	s.Run("should logout successfully", func() {
		s.T().Parallel()

		// create an organization and a user
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserPassword(validPassword))

		// create a login request
		form := url.Values{}
		form.Add("email", u.Email)
		form.Add("password", validPassword)
		loginReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf(loginPathFormat, o.Subdomain),
			strings.NewReader(form.Encode()))
		s.Require().NoError(err)
		loginReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// login
		loginRR := httptest.NewRecorder()
		h := web.SetupRoutes(s.DB, s.Config)
		h.ServeHTTP(loginRR, loginReq)

		// assert the login response
		s.Require().Equal(http.StatusOK, loginRR.Code)
		s.Empty(loginRR.Body.String())
		s.Contains(loginRR.Header().Get("Set-Cookie"), auth.JWTCookieName)

		// create a logout request
		logoutReq, err := http.NewRequest(http.MethodPost, logoutPath, nil)
		s.Require().NoError(err)
		logoutReq.Header.Add("Cookie", loginRR.Header().Get("Set-Cookie"))

		// logout
		logoutRR := httptest.NewRecorder()
		h.ServeHTTP(logoutRR, logoutReq)

		// assert the logout response
		s.Require().Equal(http.StatusOK, logoutRR.Code)
		s.Empty(logoutRR.Body.String())
		s.Contains(logoutRR.Header().Get("Set-Cookie"), "jwt_session_id=; Max-Age=0;")
	})
}
