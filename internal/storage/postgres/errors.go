package postgres

import (
	"errors"
)

var (
	ErrorInternal    = errors.New("internal database error")
	ErrDoesNotExists = errors.New("user does not exists")
)
