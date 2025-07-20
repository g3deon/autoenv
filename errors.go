package autoenv

import "errors"

var (
	ErrNilInput     = errors.New("input cannot be nil")
	ErrInvalidInput = errors.New("input must be a non-nil pointer to a struct")
)
