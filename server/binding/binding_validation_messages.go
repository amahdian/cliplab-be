package binding

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/amahdian/cliplab-be/global"
	"github.com/amahdian/cliplab-be/pkg/msg"
)

func mapValidationErrorsToString(ve validator.ValidationErrors) string {
	messages := make([]string, 0)
	for _, fe := range ve {
		messages = append(messages, fmt.Sprintf("'%s': %s", fe.Field(), translateFieldError(fe)))
	}
	return strings.Join(messages, "; ")
}

func MapValidationErrorsToMessageContainer(ve validator.ValidationErrors) *msg.MessageContainer {
	mc := msg.NewMessageContainer()
	for _, fe := range ve {
		mc.AddError(global.InvalidInputMessageGroup, fmt.Sprintf("'%s': %s", fe.Field(), translateFieldError(fe)))
	}
	return mc
}

func translateFieldError(fe validator.FieldError) string {
	if translator, ok := fieldTranslators[fe.Tag()]; ok {
		return translator(fe)
	} else {
		return fe.Error() // default error
	}
}
