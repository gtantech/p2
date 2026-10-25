package models

import "time"

type Table struct {
	Rows []TableRows
}

type TableRows struct {
	Activity     string
	Dependencies []string
	Duration     time.Duration
}
