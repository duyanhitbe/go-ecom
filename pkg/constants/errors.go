package constants

import "errors"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

var (
	ErrInvalidUsernamePassword = errors.New("invalid username or password")
	ErrUserNotFound            = errors.New("user not found")
	ErrUserAlreadyExists       = errors.New("user already exists")
)
