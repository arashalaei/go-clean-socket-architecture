package entity

import "errors"

// Domain-specific errors
var (
	ErrInvalidPerson = errors.New("invalid person")
	ErrInvalidSchool = errors.New("invalid school")
	ErrInvalidClass  = errors.New("invalid class")
	ErrNotFound      = errors.New("entity not found")
)
