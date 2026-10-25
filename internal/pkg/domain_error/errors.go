package domain_error

import "errors"

var (
	ErrNotFound  = errors.New("not found")
	ErrConflict  = errors.New("resource already exists")
	ErrReference = errors.New("resource has invalid or dependent references")
	ErrInvalid   = errors.New("invalid data")
)
