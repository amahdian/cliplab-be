package binding

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
)

type EnumType interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// EnumBinding is a binding validator that will validate if the value is one of the possible enumeration values.
// Sample usage:
//
//	type UpdateSequence struct {
//	    RigMovementType      model.SequenceRigMovementType   `json:"rigMovementType"  binding:"rig_movement_type"`
//	    RigMovementTypes     []model.SequenceRigMovementType `json:"rigMovementTypes" binding:"rig_movement_type"`
//	}
type EnumBinding[T EnumType] struct {
	tag        string
	enumValues []T
	enumStrs   []string
}

func NewEnumBinding[T EnumType](tag string, enumValues []T) *EnumBinding[T] {
	return &EnumBinding[T]{
		tag:        tag,
		enumValues: enumValues,
		enumStrs: lo.Map(enumValues, func(item T, index int) string {
			return fmt.Sprint(item)
		}),
	}
}

func (b *EnumBinding[T]) Tag() string {
	return b.tag
}

func (b *EnumBinding[T]) Translate(fe validator.FieldError) string {
	if val, ok := reflect.ValueOf(fe.Value()).Interface().(T); ok {
		return fmt.Sprintf("received value '%v' but expected one of: %q", val, strings.Join(b.enumStrs, ","))
	}
	if val, ok := reflect.ValueOf(fe.Value()).Interface().([]T); ok {
		invalidValues := lo.Without(val, b.enumValues...)
		return fmt.Sprintf("received values '%v' but expected one of: %q", invalidValues, strings.Join(b.enumStrs, ","))
	}
	return fmt.Sprintf("Unknown error while validating enum %s", b.tag)
}

func (b *EnumBinding[T]) Validate(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(T); ok {
		return lo.Contains(b.enumValues, val)
	}
	if val, ok := fl.Field().Interface().([]T); ok {
		return lo.Every(b.enumValues, val)
	}
	return false
}
