package service

import "errors"

var (
	// ErrUserWithEmailExists is returned when a user with the given email already exists.
	ErrUserWithEmailExists = errors.New("email id registered with another user")
	
	// ErrInvalidEmail is returned when the email is invalid.
	ErrInvalidEmail = errors.New("invalid email")
)
