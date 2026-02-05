package model

import "errors"

var (
	ErrInvalidURL      = errors.New("invalid URL Format")
	ErrCustomCodeTaken = errors.New("custom code already taken")
	ErrURLTooLong      = errors.New("url exceeds maximum length")
)
