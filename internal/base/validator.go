package base

import (
	"reflect"
	"strings"
	"sync"

	"github.com/camelhr/log"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

//nolint:gochecknoglobals // global variable is used to initialize the validator once
var (
	once sync.Once
	v    *validator.Validate
	uni  *ut.UniversalTranslator
)

// Validator returns a new instance of the validator. It is thread-safe.
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

		enTranslator := en.New()
		uni = ut.New(enTranslator, enTranslator)
		trans, _ := uni.GetTranslator("en")

		if err := entranslations.RegisterDefaultTranslations(v, trans); err != nil {
			log.Error("failed to register default translations: %v", err)
		}
	})
}

func ValidationTranslator() ut.Translator {
	initValidator()

	lang := "en"
	trans, _ := uni.GetTranslator(lang)

	return trans
}
