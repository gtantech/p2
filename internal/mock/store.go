package mock

import "github.com/gtantech/p2/internal/models"

type Store struct {
}

// GetHome implements [models.Store].
func (ms *Store) GetHome() models.Home {
	return models.Home{Name: "World"}
}

var _ models.Store = (*Store)(nil) //ensures MockStore implements models.Store at compile time
