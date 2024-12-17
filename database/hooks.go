package database

import (
	"context"
	"fmt"
)

// blocking
type Hooks interface {
	// execute before do query
	Before(ctx context.Context, query string, args ...interface{}) (context.Context, error)
	// execute after do query
	After(ctx context.Context, err error, query string, args ...interface{}) (context.Context, error)
}

func Compose(hooks ...Hooks) Hooks {
	return composed(hooks)
}

type composed []Hooks

func (c composed) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	var errors []error
	for _, hook := range c {
		c, err := hook.Before(ctx, query, args...)
		if err != nil {
			errors = append(errors, err)
		}
		if c != nil {
			ctx = c
		}
	}
	return ctx, wrapErrors(nil, errors)
}

func (c composed) After(ctx context.Context, err error, query string, args ...interface{}) (context.Context, error) {
	var errors []error
	for _, hook := range c {
		var e error
		c, e := hook.After(ctx, err, query, args...)
		if e != nil {
			errors = append(errors, e)
		}
		if c != nil {
			ctx = c
		}
	}
	return ctx, wrapErrors(nil, errors)
}

func wrapErrors(def error, errors []error) error {
	switch len(errors) {
	case 0:
		return def
	case 1:
		return errors[0]
	default:
		return MultipleErrors(errors)
	}
}

// MultipleErrors is an error that contains multiple errors.
type MultipleErrors []error

func (m MultipleErrors) Error() string {
	return fmt.Sprint("multiple errors:", []error(m))
}
