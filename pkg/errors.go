package pkg

import "errors"

var (
	ErrNotFound  = errors.New("Not found")
	ErrForbidden = errors.New("Forbidden")
)
