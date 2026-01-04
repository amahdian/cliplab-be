package binding

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type FieldErrorTranslator func(fe validator.FieldError) string

var fieldTranslators = make(map[string]FieldErrorTranslator) // key is binding tag

func init() {
	registerFieldTranslator("required", func(fe validator.FieldError) string {
		return fmt.Sprintf("field %q is required", fe.Field())
	})
	registerFieldTranslator("min", func(fe validator.FieldError) string {
		return fmt.Sprintf("min expected value is %q but got '%v'", fe.Param(), fe.Value())
	})
	registerFieldTranslator("max", func(fe validator.FieldError) string {
		return fmt.Sprintf("max expected value is %q but got '%v'", fe.Param(), fe.Value())
	})
	registerFieldTranslator("oneof", func(fe validator.FieldError) string {
		return fmt.Sprintf("received value '%v' but expected one of: %q", fe.Value(), fe.Param())
	})
}

func registerFieldTranslator(tag string, translator FieldErrorTranslator) {
	fieldTranslators[tag] = translator
}
