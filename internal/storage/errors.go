package storage

import (
	"errors"
)

var (
	ErrInternal      = errors.New("internal database error")
	ErrDoesNotExists = errors.New("object does not exists")
	ErrAlreadyExists = errors.New("object already exists")
	ErrUnauthorized  = errors.New("unauthorized")
)
