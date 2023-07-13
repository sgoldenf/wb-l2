package repository

import "errors"

// Repository errors
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEventNotFound = errors.New("event not found")
	ErrDuplicateID   = errors.New("duplicate event id")
)
