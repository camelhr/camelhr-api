package base

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

type (
	InvalidValidationErr = validator.InvalidValidationError
)

var (
	once sync.Once           //nolint:gochecknoglobals // global variable is used to initialize the validator once
	v    *validator.Validate //nolint:gochecknoglobals // global variable is used to initialize the validator once
)

func Validator() *validator.Validate {
	initValidator()
	return v
}

func initValidator() {
	once.Do(func() {
		v = validator.New()

		const substringPosition = 2

		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", substringPosition)[0]
			if name == "-" {
				return ""
			}

			return name
		})
	})
}
