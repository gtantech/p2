package store

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/gtantech/p2/internal/models"
)

func createTestProject(t *testing.T, store models.Model) models.Project {
	t.Helper()
	projectName := "test project"
	p, err := store.CreateProject(projectName)
	if err != nil {
		t.Fatalf("Failed to create new project. Encountered error: %v", err)
	}
	if got, want := p.ProjectDisplayName, projectName; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	return p
}

func createTableRow(project models.Project, activity models.Activity, sortRowRank int64, t *testing.T, store models.Model) models.TableRow {
	t.Helper()
	tr, err := store.CreateTableRow(project.ProjectId, activity.ActivityId, sortRowRank)
	if err != nil {
		t.Fatalf("Failed to create new table row. Encountered error: %v", err)
	}
	wantTr := models.TableRow{Activity: &activity, Dependencies: []string{}, SortRank: sortRowRank}
	if got, want := *tr.Activity, *wantTr.Activity; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := tr.Dependencies, wantTr.Dependencies; !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := tr.SortRank, wantTr.SortRank; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	return tr
}

func createActivity(activityName string, activityDuration time.Duration, p models.Project, t *testing.T, store models.Model) models.Activity {
	t.Helper()
	a, err := store.CreateActivity(p.ProjectId, activityName, activityDuration)
	if err != nil {
		t.Fatalf("Failed to create new activity. Encountered error: %v", err)
	}
	if got, want := a.DisplayName, activityName; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := a.Duration, activityDuration; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	return a
}

func testModelContract(t *testing.T, newStore func() models.Model) {
	t.Helper()

	t.Run("Create and Get Activities", func(t *testing.T) {
		s := newStore()
		p := createTestProject(t, s)

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
		a := createActivity(activityName, activityDuration, p, t, s)

		// check activity was added
		a_get, err := s.GetActivities(p.ProjectId)
		if err != nil {
			t.Fatalf("Failed to get activities. Encountered error: %v", err)
		}
		if got, want := a_get, []models.Activity{{ActivityId: a.ActivityId, DisplayName: activityName, Duration: activityDuration}}; !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Create and Get Sorted Table Row", func(t *testing.T) {
		s := newStore()
		p := createTestProject(t, s)

		// add activity
		activities := make([]models.Activity, 3)
		for i := range len(activities) {
			activityName := fmt.Sprintf("test activity %v", i+1)
			activityDuration := 5 * time.Minute
			activities[i] = createActivity(activityName, activityDuration, p, t, s)
		}
		//add table rows in reverse order (from bottom to top relative to activities)
		for i := range len(activities) {
			ib := len(activities) - 1 - i
			activityToAdd := activities[ib]
			createTableRow(p, activityToAdd, int64(ib), t, s)
		}

		//check if table rows returned is sorted correctly despite adding in reverse order
		trs, err := s.GetTableRowsSortedByRank(p.ProjectId)
		if err != nil {
			t.Fatalf("Failed to get table rows sorted by rank. Encountered error: %v", err)
		}
		for i, tr := range trs {
			if tr.ActivityId != activities[i].ActivityId {
				t.Errorf("got %v, want %v", tr.ActivityId, activities[i].ActivityId)
			}
		}
	})

	t.Run("Create Table Row At Bottom", func(t *testing.T) {
		s := newStore()
		p := createTestProject(t, s)

		// add activity
		activities := make([]models.Activity, 4)
		for i := range len(activities) {
			activityName := fmt.Sprintf("test activity %v", i+1)
			activityDuration := 5 * time.Minute
			activities[i] = createActivity(activityName, activityDuration, p, t, s)
		}
		//add table rows excluding the last activity
		for i := range len(activities) - 1 {
			activityToAdd := activities[i]
			createTableRow(p, activityToAdd, int64(i), t, s)
		}

		//add table row to the bottom
		lastActivity := activities[len(activities)-1]
		_, err := s.CreateTableRowAtBottom(p.ProjectId, lastActivity.ActivityId)
		if err != nil {
			t.Fatalf("Failed to create table row at bottom. Encountered error: %v", err)
		}

		//check if table rows returned is sorted correctly despite adding in reverse order
		trs, err := s.GetTableRowsSortedByRank(p.ProjectId)
		if err != nil {
			t.Fatalf("Failed to get table rows sorted by rank. Encountered error: %v", err)
		}
		for i, tr := range trs {
			if tr.ActivityId != activities[i].ActivityId {
				t.Errorf("got %v, want %v", tr.ActivityId, activities[i].ActivityId)
			}
		}
	})

	t.Run("Create Table Row Below", func(t *testing.T) {
		s := newStore()
		p := createTestProject(t, s)

		// add activity
		activities := make([]models.Activity, 4)
		for i := range len(activities) {
			activityName := fmt.Sprintf("test activity %v", i+1)
			activityDuration := 5 * time.Minute
			activities[i] = createActivity(activityName, activityDuration, p, t, s)
		}
		//add table rows excluding the last 2 activities
		for i := range len(activities) - 2 {
			activityToAdd := activities[i]
			createTableRow(p, activityToAdd, int64(i*1000), t, s)
		}

		//add table row for the last activity
		createTableRow(p, activities[len(activities)-1], int64(len(activities)-1)*1000, t, s)

		//add table row between the 3rd last and last activity
		_, err := s.CreateTableRowBelow(p.ProjectId, activities[len(activities)-2].ActivityId, activities[len(activities)-3].ActivityId)
		if err != nil {
			t.Fatalf("Failed to create table row below. Encountered error: %v", err)
		}

		//check if table rows returned is sorted correctly
		trs, err := s.GetTableRowsSortedByRank(p.ProjectId)
		if err != nil {
			t.Fatalf("Failed to get table rows sorted by rank. Encountered error: %v", err)
		}
		for i, tr := range trs {
			if tr.ActivityId != activities[i].ActivityId {
				t.Errorf("got %v, want %v", tr.ActivityId, activities[i].ActivityId)
			}
		}
	})
}

func TestMockStoreContract(t *testing.T) {
	testModelContract(t, func() models.Model {
		return NewMockStore()
	})
}
