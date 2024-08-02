package models

import "errors"

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrUserNotFound       = errors.New("models: user not found")
	ErrInvalidPassword    = errors.New("models: invalid password")
	ErrPermissionNotFound = errors.New("permission not found")
)
