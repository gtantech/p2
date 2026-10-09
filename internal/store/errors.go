package store

import "errors"

var (
	ErrActivityNotFound = errors.New("activity not found")
	ErrTableRowNotFound = errors.New("table row not found")
	ErrEndOfTable       = errors.New("end of table")
	ErrProjectNotFound  = errors.New("project not found")
)
