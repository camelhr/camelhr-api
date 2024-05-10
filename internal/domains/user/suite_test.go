package user_test

import (
	"testing"

	"github.com/camelhr/camelhr-api/internal/tests"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	tests.IntegrationBaseSuite
}

func TestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UserTestSuite))
}
