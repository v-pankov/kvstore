package entity

import "errors"

var (
	// ErrNotFound is generic NOTFOUND errors whose meaning depends from the context.
	ErrNotFound = errors.New("not found")
)
