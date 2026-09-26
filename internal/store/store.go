package store

import "github.com/gtantech/p2/internal/db"

type Store struct {
	Activity   *activityDbStore
	Project    *projectDbStore
	Dependency *dependencyDbStore
}

func NewStoreFromDb(queries *db.Queries) *Store {
	return &Store{
		Activity:   NewActivityStoreFromDb(queries),
		Project:    NewProjectStoreFromDb(queries),
		Dependency: NewDependencyStoreFromDb(queries),
	}
}
