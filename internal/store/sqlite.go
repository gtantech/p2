package store

import (
	"github.com/gtantech/p2/internal/database/sqlitedb"
)

func NewSqliteStore(queries sqlitedb.Queries) *SqliteStore {
	return &SqliteStore{queries: queries}
}

type SqliteStore struct {
	queries sqlitedb.Queries
}
