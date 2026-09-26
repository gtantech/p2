package routes

import "github.com/gtantech/p2/internal/store"

type storeAdapter struct {
	store *store.Store
}

// Activity implements [StoreService].
func (s *storeAdapter) Activity() ActivityStoreService {
	return s.store.Activity
}

// Dependency implements [StoreService].
func (s *storeAdapter) Dependency() DependencyStoreService {
	return s.store.Dependency
}

// Project implements [StoreService].
func (s *storeAdapter) Project() ProjectStoreService {
	return s.store.Project
}

var _ StoreService = (*storeAdapter)(nil) //ensures ExampleStruct implements ExampleInterface at compile time

func NewStoreAdapter(store *store.Store) *storeAdapter {
	return &storeAdapter{store: store}
}
