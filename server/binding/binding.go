package binding

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin/binding"
)

func Init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Ensure the Field() would return the json, uri or form tag instead of the Field name.
		// This increases the readability for the frontend
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			tagValue := ""
			if field.Tag.Get("json") != "" {
				tagValue = field.Tag.Get("json")
			} else if field.Tag.Get("uri") != "" {
				tagValue = field.Tag.Get("uri")
			} else if field.Tag.Get("form") != "" {
				tagValue = field.Tag.Get("form")
			}
			name := strings.SplitN(tagValue, ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		registerCustomBinding(v)
	}
}
