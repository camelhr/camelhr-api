package auth_test

import (
	"testing"

	"github.com/camelhr/camelhr-api/internal/tests"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	tests.IntegrationBaseSuite
}

func TestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(AuthTestSuite))
}
