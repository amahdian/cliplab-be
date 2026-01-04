package errs

import "net/http"

// An ErrorCode describes the error's category.
type ErrorCode int

const (
	// OK Returned by the Code function on a nil error. It is not a valid
	// code for an error.
	OK ErrorCode = 0

	// Unknown The error could not be categorized.
	Unknown ErrorCode = 1

	// NotFound The resource was not found.
	NotFound ErrorCode = 2

	// AlreadyExists The resource exists, but it should not.
	AlreadyExists ErrorCode = 3

	// InvalidArgument A value given to a Go CDK API is incorrect.
	InvalidArgument ErrorCode = 4

	// Internal Something unexpected happened. Internal errors always indicate
	// bugs in the Go CDK (or possibly the underlying service).
	Internal ErrorCode = 5

	// Unimplemented The feature is not implemented.
	Unimplemented ErrorCode = 6

	// FailedPrecondition The system was in the wrong state.
	FailedPrecondition ErrorCode = 7

	// PermissionDenied The caller does not have permission to execute the specified operation.
	PermissionDenied ErrorCode = 8

	// ResourceExhausted Some resource has been exhausted, typically because a service resource limit
	// has been reached.
	ResourceExhausted ErrorCode = 9

	// Canceled The operation was canceled.
	Canceled ErrorCode = 10

	// DeadlineExceeded The operation timed out.
	DeadlineExceeded ErrorCode = 11

	// Unauthenticated The authentication operation failed.
	Unauthenticated ErrorCode = 12

	// DependencyConflict The operation failed because it depended on another operations or resources.
	DependencyConflict ErrorCode = 13

	// Unavailable The service is currently unavailable.
	Unavailable ErrorCode = 14

	// ResourceLocked The resource is locked.
	ResourceLocked ErrorCode = 15

	// ParseError The expression parsing failed.
	ParseError ErrorCode = 16

	// CompilationError The expression compilation failed.
	CompilationError ErrorCode = 17

	// RuntimeEvaluationError The expression evaluation failed.
	RuntimeEvaluationError ErrorCode = 18

	// TypeCheckError The expression result type checking failed.
	TypeCheckError ErrorCode = 19
)

func (i ErrorCode) String() string {
	switch i {
	case NotFound:
		return "NotFound"
	case AlreadyExists:
		return "AlreadyExists"
	case InvalidArgument:
		return "InvalidArgument"
	case Internal:
		return "Internal"
	case Unimplemented:
		return "Unimplemented"
	case FailedPrecondition:
		return "FailedPrecondition"
	case PermissionDenied:
		return "PermissionDenied"
	case ResourceExhausted:
		return "ResourceExhausted"
	case Canceled:
		return "Canceled"
	case DeadlineExceeded:
		return "DeadlineExceeded"
	case Unauthenticated:
		return "Unauthenticated"
	case DependencyConflict:
		return "DependencyConflict"
	case Unavailable:
		return "Unavailable"
	case ResourceLocked:
		return "ResourceLocked"
	case ParseError:
		return "ParseError"
	case CompilationError:
		return "CompilationError"
	case RuntimeEvaluationError:
		return "RuntimeEvaluationError"
	case TypeCheckError:
		return "TypeCheckError"
	}
	return "Unknown"
}

func (i ErrorCode) HttpStatus() int {
	switch i {
	case NotFound:
		return http.StatusNotFound
	case AlreadyExists:
		return http.StatusBadRequest
	case InvalidArgument:
		return http.StatusBadRequest
	case Internal:
		return http.StatusInternalServerError
	case Unimplemented:
		return http.StatusNotImplemented
	case FailedPrecondition:
		return http.StatusPreconditionFailed
	case PermissionDenied:
		return http.StatusForbidden
	case ResourceExhausted:
		return http.StatusTooManyRequests
	case Canceled:
		return http.StatusServiceUnavailable
	case DeadlineExceeded:
		return http.StatusGatewayTimeout
	case Unauthenticated:
		return http.StatusUnauthorized
	case DependencyConflict:
		return http.StatusConflict
	case Unavailable:
		return http.StatusServiceUnavailable
	case ResourceLocked:
		return http.StatusLocked
	case ParseError:
		return http.StatusBadRequest
	case CompilationError:
		return http.StatusInternalServerError
	case RuntimeEvaluationError:
		return http.StatusInternalServerError
	case TypeCheckError:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func (i ErrorCode) MessageGroup() string {
	switch i {
	case NotFound:
		return "Not Found"
	case AlreadyExists:
		return "Already Exists"
	case InvalidArgument:
		return "Invalid Input"
	case Internal:
		return "Internal Issue"
	case Unimplemented:
		return "Unimplemented"
	case FailedPrecondition:
		return "Failed Precondition"
	case PermissionDenied:
		return "Permission Denied"
	case ResourceExhausted:
		return "Resource Exhausted"
	case Canceled:
		return "Canceled"
	case DeadlineExceeded:
		return "Deadline Exceeded"
	case Unauthenticated:
		return "Unauthenticated"
	case DependencyConflict:
		return "Dependency Conflict"
	case Unavailable:
		return "Unavailable"
	case ResourceLocked:
		return "Resource Locked"
	case ParseError:
		return "Parse Error"
	case CompilationError:
		return "Compilation Error"
	case RuntimeEvaluationError:
		return "Runtime Evaluation Error"
	case TypeCheckError:
		return "Type Check Error"
	}
	return "Unknown"
}
