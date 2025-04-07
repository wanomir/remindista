package domain

import "errors"

var (
	ErrorInvalidCallback = errors.New("invalid callback")
	ErrorShortTag        = errors.New("tag should be at least 2 characters long")
)
