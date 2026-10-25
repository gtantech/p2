package store

import "github.com/gtantech/p2/internal/models"

type MockStore struct {
}

// GetHome implements [models.Store].
func (ms *MockStore) GetHome() models.Home {
	return models.Home{Name: "World"}
}

var _ models.Store = (*MockStore)(nil) //ensures MockStore implements models.Store at compile time
