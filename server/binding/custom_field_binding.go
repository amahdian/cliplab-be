package binding

import (
	"github.com/go-playground/validator/v10"
)

type CustomFieldBinding interface {
	Tag() string
	Translate(fe validator.FieldError) string
	Validate(fl validator.FieldLevel) bool
}
