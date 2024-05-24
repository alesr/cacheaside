package repository

import "errors"

// Enumerate the errors that can be returned by the repository.

var (
	ErrItemNotFound = errors.New("item not found")
	ErrInvalidItem  = errors.New("invalid item")
)
