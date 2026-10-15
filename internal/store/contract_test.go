package store

import (
	"slices"
	"testing"
	"time"

	"github.com/gtantech/p2/internal/models"
)

func testModelContract(t *testing.T, newStore func() models.Model) {
	t.Helper()

	t.Run("Create and Get Activities", func(t *testing.T) {
		s := newStore()
		projectName := "test project"
		p, err := s.CreateProject(projectName)
		if err != nil {
			t.Fatalf("Failed to create new project. Encountered error: %v", err)
		}
		if got, want := p.ProjectDisplayName, projectName; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}

		// check no activites present
		a_getEmpty, err := s.GetActivities(p.ProjectId)
		if err != nil {
			t.Fatalf("Failed to get activities. Encountered error: %v", err)
		}
		if got, want := a_getEmpty, []models.Activity{}; !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}

		// add activity
		activityName := "test activity"
		activityDuration := 5 * time.Minute
		a, err := s.CreateActivity(p.ProjectId, activityName, activityDuration)
		if err != nil {
			t.Fatalf("Failed to create new activity. Encountered error: %v", err)
		}
		if got, want := a.DisplayName, activityName; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
		if got, want := a.Duration, activityDuration; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}

		// check activity was added
		a_get, err := s.GetActivities(p.ProjectId)
		if err != nil {
			t.Fatalf("Failed to get activities. Encountered error: %v", err)
		}
		if got, want := a_get, []models.Activity{{ActivityId: a.ActivityId, DisplayName: activityName, Duration: activityDuration}}; !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestMockStoreContract(t *testing.T) {
	testModelContract(t, func() models.Model {
		return NewMockStore()
	})
}
