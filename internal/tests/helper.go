package tests

import (
	"strings"

	"github.com/stretchr/testify/mock"
)

func QueryMatcher(queryLabel string) any {
	return mock.MatchedBy(func(a any) bool {
		if query, ok := a.(string); ok {
			return strings.Contains(query, queryLabel)
		}

		return false
	})
}
