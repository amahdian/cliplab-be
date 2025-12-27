package global

type MessageGroup string

const (
	//Internal
	UnknownMessageGroup       = "unknown"
	InternalIssueMessageGroup = "Internal Issue"

	//Permission
	PermissionDeniedMessageGroup = "Permission denied"

	//Validation
	InvalidInputMessageGroup = "Invalid input"
	NotFoundMessageGroup     = "Not found"
)
