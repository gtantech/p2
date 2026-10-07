package store

import (
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models"
)

type MockStore struct {
}

// GetDependencies implements [models.Model].
func (ms *MockStore) GetDependencies(activityId uuid.UUID) ([]models.Activity, error) {

	activities := []models.Activity{
		{ActivityId: uuid.MustParse("dcce6c98-b31b-4ed6-a2ed-15d24ae96b41"), DisplayName: "A", Duration: 5 * time.Minute},
		{ActivityId: uuid.MustParse("a4518b28-e597-4295-a532-3a501a75ab2b"), DisplayName: "B", Duration: 4 * time.Minute},
		{ActivityId: uuid.MustParse("aa5a833c-2407-464c-93c4-fa9edba03b52"), DisplayName: "C", Duration: 3 * time.Minute},
	}

	dependencies := map[uuid.UUID][]models.Activity{}
	dependencies[uuid.MustParse("a4518b28-e597-4295-a532-3a501a75ab2b")] = []models.Activity{activities[0]}
	dependencies[uuid.MustParse("aa5a833c-2407-464c-93c4-fa9edba03b52")] = []models.Activity{activities[0], activities[1]}
	return dependencies[activityId], nil
}

// GetActivities implements [models.Model].
func (ms *MockStore) GetActivities(projectId uuid.UUID) ([]models.Activity, error) {
	activities := []models.Activity{
		{ActivityId: uuid.MustParse("dcce6c98-b31b-4ed6-a2ed-15d24ae96b41"), DisplayName: "A", Duration: 5 * time.Minute},
		{ActivityId: uuid.MustParse("a4518b28-e597-4295-a532-3a501a75ab2b"), DisplayName: "B", Duration: 4 * time.Minute},
		{ActivityId: uuid.MustParse("aa5a833c-2407-464c-93c4-fa9edba03b52"), DisplayName: "C", Duration: 3 * time.Minute},
	}
	return activities, nil
}

// GetHome implements [models.Model].
func (ms *MockStore) GetHome() models.Home {
	return models.Home{Name: "World"}
}

var _ models.Model = (*MockStore)(nil) //ensures MockStore implements models.Model at compile time
