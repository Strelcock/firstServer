package auth

import "errors"

var (
	ErrUserExists       = errors.New("user exists")
	ErrWrongCredentials = errors.New("wrong email or password")
)
