// Package errs provides a wrapper around github/go/errors that adds the
// ability to associate a kvp type context with the error.
package errs

import (
	"github.com/github/go/errors"
	"github.com/github/mu/kvp"
)

// New returns an error with the supplied message.
func New(message string) error {
	return errors.New(message)
}

// Errorf formats according to a format specifier and returns the string as a
// value that satisfies error.
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// With creates a new error with an associated kvp context.
func With(message string, fields ...kvp.Field) error {
	return WrapWith(New(message), fields...)
}

// Wrap returns an error annotating err with message. If err is nil, Wrap
// returns nil.
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf returns an error annotating err with args. If err is nil, Wrapf
// returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// WrapWith returns an error with an associated kvp context.
func WrapWith(err error, fields ...kvp.Field) error {
	return _error{
		err: errors.Wrap(err, err.Error()),
		ctx: kvp.New(fields...),
	}
}

// Context returns the associated kvp context, if possible.
func Context(err error) *kvp.KVP {
	type contexter interface {
		Context() *kvp.KVP
	}

	ctxer, ok := err.(contexter)
	if !ok {
		return nil
	}

	return ctxer.Context()
}

// _error is an error implementation returned by the methods that add a kvp
// context.
type _error struct {
	err error
	ctx *kvp.KVP
}

// Error returns the underlying error's Error().
func (e _error) Error() string {
	return e.err.Error()
}

// Context returns the kvp context.
func (e _error) Context() *kvp.KVP {
	return e.ctx
}

// StackTrace allows the wrapped errors to use the underlying errors package's
// stack trace capabilities. When returning the stacktrace, we chop off the top
// two frames, which will be errs package constructors.
func (e _error) StackTrace() errors.StackTrace {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if err, ok := e.err.(stackTracer); ok {
		st := err.StackTrace()
		return st[2 : len(st)-1]
	}

	return nil
}
