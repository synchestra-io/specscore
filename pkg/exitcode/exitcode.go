// Package exitcode defines the shared exit code constants and error type
// used by all SpecScore CLI commands and library consumers.
package exitcode

import "fmt"

// Standard exit codes shared by every CLI command.
const (
	Success      = 0  // Operation completed successfully.
	Conflict     = 1  // Concurrent-modification conflict.
	InvalidArgs  = 2  // Missing or invalid command arguments/flags.
	NotFound     = 3  // Requested resource does not exist.
	InvalidState = 4  // State transition is not allowed.
	Unexpected   = 10 // Catch-all for unexpected runtime errors.
)

// Error carries a machine-readable exit code alongside a human-readable
// message. It satisfies both the error interface and the ExitCode()
// convention checked by the top-level CLI runner.
type Error struct {
	code int
	msg  string
}

func (e *Error) Error() string { return e.msg }

// ExitCode returns the numeric exit code for this error.
func (e *Error) ExitCode() int { return e.code }

// New creates an Error with the given exit code and message.
func New(code int, msg string) *Error {
	return &Error{code: code, msg: msg}
}

// Newf creates an Error with the given exit code and formatted message.
func Newf(code int, format string, args ...any) *Error {
	return &Error{code: code, msg: fmt.Sprintf(format, args...)}
}

// --- Convenience constructors for standard exit codes ---

// ConflictError returns an exit-code-1 error.
func ConflictError(msg string) *Error { return &Error{code: Conflict, msg: msg} }

// ConflictErrorf returns an exit-code-1 error with a formatted message.
func ConflictErrorf(format string, args ...any) *Error {
	return Newf(Conflict, format, args...)
}

// InvalidArgsError returns an exit-code-2 error.
func InvalidArgsError(msg string) *Error { return &Error{code: InvalidArgs, msg: msg} }

// InvalidArgsErrorf returns an exit-code-2 error with a formatted message.
func InvalidArgsErrorf(format string, args ...any) *Error {
	return Newf(InvalidArgs, format, args...)
}

// NotFoundError returns an exit-code-3 error.
func NotFoundError(msg string) *Error { return &Error{code: NotFound, msg: msg} }

// NotFoundErrorf returns an exit-code-3 error with a formatted message.
func NotFoundErrorf(format string, args ...any) *Error {
	return Newf(NotFound, format, args...)
}

// InvalidStateError returns an exit-code-4 error.
func InvalidStateError(msg string) *Error { return &Error{code: InvalidState, msg: msg} }

// InvalidStateErrorf returns an exit-code-4 error with a formatted message.
func InvalidStateErrorf(format string, args ...any) *Error {
	return Newf(InvalidState, format, args...)
}

// UnexpectedError returns an exit-code-10 error.
func UnexpectedError(msg string) *Error { return &Error{code: Unexpected, msg: msg} }

// UnexpectedErrorf returns an exit-code-10 error with a formatted message.
func UnexpectedErrorf(format string, args ...any) *Error {
	return Newf(Unexpected, format, args...)
}
