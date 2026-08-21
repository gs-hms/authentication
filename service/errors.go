package service

import "errors"

var (
	// ErrUserWithEmailExists is returned when a user with the given email already exists.
	ErrUserWithEmailExists = errors.New("email id registered with another user")

	// ErrInvalidEmail is returned when the email is invalid.
	ErrInvalidEmail = errors.New("invalid email")

	// ErrUserWithPhoneExists is returned when a user with the given phone number already exists.
	ErrUserWithPhoneExists = errors.New("phone number registered with another user")

	// ErrInvalidPhone is returned when the phone number is invalid.
	ErrInvalidPhone = errors.New("invalid phone number")
)
