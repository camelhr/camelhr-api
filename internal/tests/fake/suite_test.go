package fake_test

import (
	"testing"

	"github.com/camelhr/camelhr-api/internal/tests"
	"github.com/stretchr/testify/suite"
)

type FakeTestSuite struct {
	tests.BaseSuite
}

func TestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(FakeTestSuite))
}
