package errs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/xerrors"
)

type Error struct {
	Code   ErrorCode
	msg    string
	frame  xerrors.Frame
	err    error
	format string
	args   []interface{}
}

func (e *Error) Error() string {
	return fmt.Sprint(e)
}

func (e *Error) FormatError(p xerrors.Printer) (next error) {
	if e.msg == "" {
		p.Printf("Code: %v", e.Code)
	} else {
		p.Printf("%s", e.msg)
	}
	e.frame.Format(p)
	return e.err
}

func (e *Error) Format(s fmt.State, c rune) {
	xerrors.FormatError(e, s, c)
}

// Unwrap returns the error underlying the receiver, which may be nil.
func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) DeepCopy() *Error {
	clonedArgs := e.args
	return &Error{
		Code:   e.Code,
		msg:    e.msg,
		frame:  e.frame,
		err:    e.err,
		format: e.format,
		args:   clonedArgs,
	}
}

func new(c ErrorCode, err error, callDepth int, msg string, format string, args []interface{}) *Error {
	return &Error{
		Code:   c,
		msg:    msg,
		frame:  xerrors.Caller(callDepth),
		err:    err,
		format: format,
		args:   args,
	}
}

// New returns a new error with the given code, underlying error and message. Pass 1
// for the call depth if New is called from the function raising the error; pass 2 if
// it is called from a helper function that was invoked by the original function; and
// so on.
func New(c ErrorCode, err error, callDepth int, msg string) *Error {
	return new(c, err, callDepth, msg, msg, make([]interface{}, 0))
}

// Newf uses format and args to format a message, then calls New.
func Newf(c ErrorCode, err error, format string, args ...any) *Error {
	return new(c, err, 2, fmt.Sprintf(format, args...), format, args)
}

// Wrapf detect the underlying error code, uses format and args to format a message, then calls New.
func Wrapf(err error, format string, args ...any) *Error {
	return new(Code(err), err, 2, fmt.Sprintf(format, args...), format, args)
}

func Code(err error) ErrorCode {
	if err == nil {
		return OK
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	if errors.Is(err, context.Canceled) {
		return Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return DeadlineExceeded
	}
	return Unknown
}

func IsCode(err error, code ErrorCode) bool {
	return Code(err) == code
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, &target)
}

func Message(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) {
		return e.msg
	}
	return ""
}

func Format(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) {
		return e.format
	}
	return ""
}

func Args(err error) []interface{} {
	if err == nil {
		return make([]interface{}, 0)
	}
	var e *Error
	if errors.As(err, &e) {
		return e.args
	}
	return make([]interface{}, 0)
}

type CompositeErr struct {
	errs []error
}

func NewCompositeErr() *CompositeErr {
	return &CompositeErr{errs: make([]error, 0)}
}

func (c *CompositeErr) Error() string {
	n := len(c.errs)
	if n == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i, e := range c.errs {
		sb.WriteString(e.Error())
		if i < n-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func (c *CompositeErr) Add(err error) {
	if err != nil {
		c.errs = append(c.errs, err)
	}
}

// DeepCopy creates a deep copy of the CompositeErr
func (c *CompositeErr) DeepCopy() *CompositeErr {
	if c == nil {
		return nil
	}

	// Create a new instance
	copyErr := &CompositeErr{
		errs: make([]error, len(c.errs)),
	}

	var ce *CompositeErr
	// Copy each error
	for i, err := range c.errs {
		if errors.As(err, &ce) {
			copyErr.errs[i] = ce.DeepCopy()
		}
	}

	var e *Error
	// Copy each error
	for i, err := range c.errs {
		if errors.As(err, &e) {
			copyErr.errs[i] = e.DeepCopy()
		}
	}

	return copyErr
}

func Errors(errs []error) error {
	return &CompositeErr{errs: errs}
}

func HasErrors(err error) bool {
	hasErrors := false

	if err == nil {
		return false
	}

	var comp *CompositeErr
	if errors.As(err, &comp) {
		if len(comp.errs) == 0 {
			return false
		}
		for _, e := range comp.errs {
			if errors.As(e, &comp) {
				hasErrors = hasErrors || HasErrors(e)
			} else if e == nil {
				return false
			} else {
				hasErrors = true
			}
		}
	}

	return hasErrors
}

func FindErrorWithCode(err error, code ErrorCode) *Error {
	if err == nil {
		return nil
	}

	var comp *CompositeErr
	if errors.As(err, &comp) {
		for _, e := range comp.errs {
			if found := FindErrorWithCode(e, code); found != nil {
				return found
			}
		}
		return nil
	}

	var e *Error
	if errors.As(err, &e) {
		if e.Code == code {
			return e
		}
		if found := FindErrorWithCode(e.Unwrap(), code); found != nil {
			return found
		}
	} else {
		if unwrapped := errors.Unwrap(err); unwrapped != nil {
			return FindErrorWithCode(unwrapped, code)
		}
	}

	return nil
}

func DeepCopy(err error) error {
	if err == nil {
		return nil
	}

	var comp *CompositeErr
	if errors.As(err, &comp) {
		return comp.DeepCopy()
	}

	var e *Error
	if errors.As(err, &e) {
		return e.DeepCopy()
	}

	return nil
}

func FindErrorWithCodeAndArg(err error, code ErrorCode, arg interface{}) *Error {
	if err == nil {
		return nil
	}

	var comp *CompositeErr
	if errors.As(err, &comp) {
		for _, e := range comp.errs {
			if found := FindErrorWithCodeAndArg(e, code, arg); found != nil {
				return found
			}
		}
		return nil
	}

	var e *Error
	if errors.As(err, &e) {
		if e.Code == code {
			containsArg := lo.Contains(e.args, arg)
			if containsArg {
				return e
			} else {
				return nil
			}
		}
		if found := FindErrorWithCodeAndArg(e.Unwrap(), code, arg); found != nil {
			return found
		}
	} else {
		if unwrapped := errors.Unwrap(err); unwrapped != nil {
			return FindErrorWithCodeAndArg(unwrapped, code, arg)
		}
	}

	return nil
}

func RemoveWithErrorCode(err error, code ErrorCode) {
	searchedError := FindErrorWithCode(err, code)
	if searchedError != nil {
		var comp *CompositeErr
		if errors.As(err, &comp) {
			for index, e := range comp.errs {
				if errors.As(e, &comp) {
					RemoveWithErrorCode(e, code)
				} else {
					if searchedError.Code == Code(e) {
						if len(comp.errs) == 1 {
							comp.errs = make([]error, 0)
						} else {
							comp.errs = append(comp.errs[:index], comp.errs[index+1:]...)
						}
						return
					}
				}
			}
		}
	}
}
