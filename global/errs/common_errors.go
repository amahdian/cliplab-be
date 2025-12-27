package errs

import (
	"fmt"
)

const (
	DefaultEntryNotFoundMessage       = "%q entry by id %s could not be found."
	InvalidMobileMessage              = "Mobile number is invalid."
	InvalidSearchFieldMessage         = "Invalid field %q for search condition."
	InvalidSequenceGroupAccessMessage = "User does not have access to sequence group %d"
	FailedToListItemsMessage          = "Failed to list %q"
	InvalidColorMessage               = "Color should be a valid hex color."
	InvalidTimeOfDayMessage           = "Invalid time input: %s. Time should be in this pattern: \"21:00\""
)

func NewInvalidSearchFieldErr(fieldName string) error {
	return fmt.Errorf(InvalidSearchFieldMessage, fieldName)
}

func NewInvalidSequenceGroupAccessErr(groupId int64) error {
	return fmt.Errorf(InvalidSequenceGroupAccessMessage, groupId)
}

type EntryNotFoundErr struct {
	message string
}

func NewEntryNotFoundErr(message string) EntryNotFoundErr {
	return EntryNotFoundErr{message}
}

func (e EntryNotFoundErr) Error() string {
	return e.message
}
