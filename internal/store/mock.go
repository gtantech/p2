package store

import (
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models"
)

type MockStore struct {
}

// GetActivities implements [models.Model].
func (ms *MockStore) GetActivities(projectId uuid.UUID) ([]models.Activity, error) {
	activities := []models.Activity{
		{ActivityId: uuid.New(), DisplayName: "A", Duration: 5 * time.Minute},
		{ActivityId: uuid.New(), DisplayName: "B", Duration: 4 * time.Minute},
		{ActivityId: uuid.New(), DisplayName: "C", Duration: 3 * time.Minute},
	}
	return activities, nil
}

// GetHome implements [models.Model].
func (ms *MockStore) GetHome() models.Home {
	return models.Home{Name: "World"}
}

var _ models.Model = (*MockStore)(nil) //ensures MockStore implements models.Model at compile time
