package store

import "github.com/gtantech/p2/internal/db"

type Store struct {
	Activity   ActivityStore
	Project    ProjectStore
	Dependency DependencyStore
}

func NewStoreFromDb(queries *db.Queries) *Store {
	return &Store{
		Activity:   NewActivityStoreFromDb(queries),
		Project:    NewProjectStoreFromDb(queries),
		Dependency: NewDependencyStoreFromDb(queries),
	}
}
