package tests

import (
	"testing"

	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/stretchr/testify/suite"
)

type IntegrationBaseSuite struct {
	suite.Suite
	DB          database.Database
	PGContainer *PostgreSQLContainer
}

func (s *IntegrationBaseSuite) SetupSuite() {
	if testing.Short() {
		s.T().Skip("skipping suite in short mode")
	}

	setupDone := false
	defer func() {
		if !setupDone {
			s.TearDownSuite()
		}
	}()

	c, err := NewPostgresContainer()
	s.Require().NoError(err)
	s.PGContainer = c

	db, err := c.Connect()
	s.Require().NoError(err)
	s.DB = database.NewPostgresDatabase(db)

	err = RunMigrations(db.DB)
	s.Require().NoError(err)

	setupDone = true
}

func (s *IntegrationBaseSuite) TearDownSuite() {
	err := s.PGContainer.Purge()
	s.Require().NoError(err)
}
