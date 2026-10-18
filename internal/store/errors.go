package store

import "errors"

var (
	ErrActivityNotFound = errors.New("activity not found")
	ErrProjectNotFound  = errors.New("project not found")
)
