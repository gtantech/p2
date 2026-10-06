package store

import "github.com/gtantech/p2/internal/models"

type MockStore struct {
}

// GetHome implements [models.Model].
func (ms *MockStore) GetHome() models.Home {
	return models.Home{Name: "World"}
}

var _ models.Model = (*MockStore)(nil) //ensures MockStore implements models.Model at compile time
