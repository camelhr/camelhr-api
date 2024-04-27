package fake_test

import (
	"testing"

	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/camelhr/camelhr-api/internal/tests"
	"github.com/stretchr/testify/suite"
)

type FakeTestSuite struct {
	suite.Suite
	db          database.Database
	pgContainer *tests.PostgreSQLContainer
}

func TestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(FakeTestSuite))
}

func (s *FakeTestSuite) SetupSuite() {
	setupDone := false
	defer func() {
		if !setupDone {
			// since teardown suite is not called when setup suite fails
			// https://github.com/stretchr/testify/issues/1123
			s.TearDownSuite()
		}
	}()

	c, err := tests.NewPostgresContainer("5433")
	s.Require().NoError(err)
	s.pgContainer = c

	db, err := c.Connect()
	s.Require().NoError(err)
	s.db = database.NewPostgresDatabase(db)

	err = tests.RunMigrations(db.DB)
	s.Require().NoError(err)

	setupDone = true
}

func (s *FakeTestSuite) TearDownSuite() {
	err := s.pgContainer.Purge()
	s.Require().NoError(err)
}
