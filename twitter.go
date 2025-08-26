package gotwitter

import "errors"

var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
	ErrBadCredentials = errors.New("email/password combination is incorrect")
	ErrInvalidAccessToken = errors.New("invalid access token")
)
