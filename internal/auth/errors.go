package auth

import "errors"

var (
	ErrExists     = errors.New("user already exists")
	ErrWrongCreds = errors.New("invalid username or password")
)
