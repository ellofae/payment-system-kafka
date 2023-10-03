package errors

import "errors"

var (
	ErrNoRecordFound       = errors.New("no such requested record has been found in database relationship")
	ErrIncorrectPassword   = errors.New("incorrect password was provided")
	ErrIncorrectEmail      = errors.New("incorrect email was provided")
	ErrAlreadyExists       = errors.New("record already exists")
	ErrGenerateAccessToken = errors.New("unable to generate an access token")
)
