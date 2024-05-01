package fake

import (
	"context"

	"github.com/stretchr/testify/mock"
)

//nolint:gochecknoglobals // global variables are used for mocking
var (
	MockBool      = mock.AnythingOfType("bool")
	MockBoolPtr   = mock.AnythingOfType("*bool")
	MockInt64     = mock.AnythingOfType("int64")
	MockInt64Ptr  = mock.AnythingOfType("*int64")
	MockString    = mock.AnythingOfType("string")
	MockStringPtr = mock.AnythingOfType("*string")
	MockTime      = mock.AnythingOfType("time.Time")
	MockTimePtr   = mock.AnythingOfType("*time.Time")
	MockContext   = mock.MatchedBy(func(_ context.Context) bool { return true })
)
