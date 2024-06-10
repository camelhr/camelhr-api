package auth_test

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/auth"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
)

func (s *AuthTestSuite) TestServiceIntegration_Register() {
	s.Run("should return error when subdomain already exists", func() {
		s.T().Parallel()

		userRepo := user.NewRepository(s.DB)
		userService := user.NewService(userRepo)
		orgRepo := organization.NewRepository(s.DB)
		orgService := organization.NewService(orgRepo)
		authService := auth.NewService(s.DB, orgService, userService, s.Config.AppSecret)

		subdomain := gofakeit.LetterN(20)
		orgName := gofakeit.LetterN(50)
		email := gofakeit.Email()
		password := "yiG3@#fj"

		err := authService.Register(context.Background(), email, password, subdomain, orgName)
		s.Require().NoError(err)

		newOrg, err := orgService.GetOrganizationBySubdomain(context.Background(), subdomain)
		s.Require().NoError(err)
		s.NotZero(newOrg.ID)
		s.Equal(subdomain, newOrg.Subdomain)
		s.Equal(orgName, newOrg.Name)

		newUser, err := userService.GetUserByOrgIDEmail(context.Background(), newOrg.ID, email)
		s.Require().NoError(err)
		s.NotZero(newUser.ID)
		s.Equal(email, newUser.Email)
		s.NotEmpty(newUser.PasswordHash)
	})

	s.Run("should rollback txn and revert org if user creation fails", func() {
		s.T().Parallel()

		userRepo := user.NewRepository(s.DB)
		userService := user.NewService(userRepo)
		orgRepo := organization.NewRepository(s.DB)
		orgService := organization.NewService(orgRepo)
		authService := auth.NewService(s.DB, orgService, userService, s.Config.AppSecret)

		subdomain := gofakeit.LetterN(20)
		orgName := gofakeit.LetterN(50)

		// provide an invalid email to trigger user creation error
		email := "@@@invalid"
		password := "niG3@#fj"

		err := authService.Register(context.Background(), email, password, subdomain, orgName)
		s.Require().Error(err)
		s.Require().ErrorContains(err, "email must be a valid email address")

		_, err = orgService.GetOrganizationBySubdomain(context.Background(), subdomain)
		s.Require().Error(err)
		s.Require().EqualError(err, "organization not found for the given subdomain")
	})

	s.Run("should register a new organization with owner", func() {
		s.T().Parallel()

		userRepo := user.NewRepository(s.DB)
		userService := user.NewService(userRepo)
		orgRepo := organization.NewRepository(s.DB)
		orgService := organization.NewService(orgRepo)
		authService := auth.NewService(s.DB, orgService, userService, s.Config.AppSecret)

		subdomain := gofakeit.LetterN(20)
		orgName := gofakeit.LetterN(50)
		email := gofakeit.Email()
		password := validPassword

		err := authService.Register(context.Background(), email, password, subdomain, orgName)
		s.Require().NoError(err)

		newOrg, err := orgService.GetOrganizationBySubdomain(context.Background(), subdomain)
		s.Require().NoError(err)
		s.NotZero(newOrg.ID)
		s.Equal(subdomain, newOrg.Subdomain)
		s.Equal(orgName, newOrg.Name)
		s.NotNil(newOrg.DisabledAt)
		s.Nil(newOrg.DeletedAt)
		s.WithinDuration(time.Now().UTC(), *newOrg.DisabledAt, 1*time.Minute)
		s.NotNil(newOrg.Comment)
		s.Equal(auth.NewOrgDisableComment, *newOrg.Comment)

		newUser, err := userService.GetUserByOrgIDEmail(context.Background(), newOrg.ID, email)
		s.Require().NoError(err)
		s.NotZero(newUser.ID)
		s.Nil(newUser.DeletedAt)
		s.Nil(newUser.DisabledAt)
		s.Equal(email, newUser.Email)
		s.NotEmpty(newUser.PasswordHash)
	})
}

func (s *AuthTestSuite) TestServiceIntegration_Login() {
	s.Run("should return error when user is disabled", func() {
		s.T().Parallel()

		userRepo := user.NewRepository(s.DB)
		userService := user.NewService(userRepo)
		orgRepo := organization.NewRepository(s.DB)
		orgService := organization.NewService(orgRepo)
		authService := auth.NewService(s.DB, orgService, userService, s.Config.AppSecret)

		password := "2iG3@#fj"
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserPassword(password))

		token, err := authService.Login(context.Background(), o.Subdomain, u.Email, password)
		s.Require().NoError(err)
		s.NotEmpty(token)
	})

	s.Run("should login successfully", func() {
		s.T().Parallel()

		userRepo := user.NewRepository(s.DB)
		userService := user.NewService(userRepo)
		orgRepo := organization.NewRepository(s.DB)
		orgService := organization.NewService(orgRepo)
		authService := auth.NewService(s.DB, orgService, userService, s.Config.AppSecret)

		password := validPassword
		o := fake.NewOrganization(s.DB)
		u := fake.NewUser(s.DB, o.ID, fake.UserPassword(password))

		token, err := authService.Login(context.Background(), o.Subdomain, u.Email, password)
		s.Require().NoError(err)
		s.NotEmpty(token)
	})
}
