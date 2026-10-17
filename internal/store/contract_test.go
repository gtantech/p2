package store

import (
	"fmt"
	"maps"
	"slices"
	"testing"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/database/sqlitedb"
	"github.com/gtantech/p2/internal/models"
)

func createTestProject(t *testing.T, store models.Model) models.Project {
	t.Helper()
	projectName := "test project"
	p, err := store.CreateProject(t.Context(), projectName)
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
	tr, err := store.CreateTableRow(t.Context(), project.ProjectId, activity.ActivityId, sortRowRank)
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
	a, err := store.CreateActivity(t.Context(), p.ProjectId, activityName, activityDuration)
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

func assertActivity(gotActivity models.Activity, wantActivity models.Activity, t *testing.T) {
	t.Helper()
	if got, want := gotActivity.ActivityId, wantActivity.ActivityId; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := gotActivity.DisplayName, wantActivity.DisplayName; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := gotActivity.Duration, wantActivity.Duration; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func testModelContract(t *testing.T, newStore func() models.Model) {
	t.Helper()

	t.Run("Create and Get Activities", func(t *testing.T) {
		s := newStore()
		p := createTestProject(t, s)

		// check no activites present
		a_getEmpty, err := s.GetActivities(t.Context(), p.ProjectId)
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
		a_get, err := s.GetActivities(t.Context(), p.ProjectId)
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
		trs, err := s.GetTableRowsSortedByRank(t.Context(), p.ProjectId)
		if err != nil {
			t.Fatalf("Failed to get table rows sorted by rank. Encountered error: %v", err)
		}
		for i, tr := range trs {
			assertActivity(*tr.Activity, activities[i], t)
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
		_, err := s.CreateTableRowAtBottom(t.Context(), p.ProjectId, lastActivity.ActivityId)
		if err != nil {
			t.Fatalf("Failed to create table row at bottom. Encountered error: %v", err)
		}

		//check if table rows returned is sorted correctly despite adding in reverse order
		trs, err := s.GetTableRowsSortedByRank(t.Context(), p.ProjectId)
		if err != nil {
			t.Fatalf("Failed to get table rows sorted by rank. Encountered error: %v", err)
		}
		for i, tr := range trs {
			assertActivity(*tr.Activity, activities[i], t)
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
		_, err := s.CreateTableRowBelow(t.Context(), p.ProjectId, activities[len(activities)-2].ActivityId, activities[len(activities)-3].ActivityId)
		if err != nil {
			t.Fatalf("Failed to create table row below. Encountered error: %v", err)
		}

		//check if table rows returned is sorted correctly
		trs, err := s.GetTableRowsSortedByRank(t.Context(), p.ProjectId)
		if err != nil {
			t.Fatalf("Failed to get table rows sorted by rank. Encountered error: %v", err)
		}
		for i, tr := range trs {
			assertActivity(*tr.Activity, activities[i], t)
		}
	})

	t.Run("Create Dependency", func(t *testing.T) {
		s := newStore()
		p := createTestProject(t, s)

		// add activity
		activities := make([]models.Activity, 3)
		for i := range len(activities) {
			activityName := fmt.Sprintf("test activity %v", i+1)
			activityDuration := 5 * time.Minute
			activities[i] = createActivity(activityName, activityDuration, p, t, s)
		}
		//add table rows
		for i := range len(activities) {
			activityToAdd := activities[i]
			createTableRow(p, activityToAdd, int64(i), t, s)
		}

		d, err := s.CreateDependency(t.Context(), p.ProjectId, activities[1].ActivityId, activities[0].ActivityId)
		if err != nil {
			t.Fatalf("Failed to create dependency. Encountered error: %v", err)
		}
		if got, want := d.SuccessorActivityId, activities[1].ActivityId; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
		if got, want := d.PredecessorActivityId, activities[0].ActivityId; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}

		wantDependencies := map[uuid.UUID][]string{}
		wantDependencies[activities[0].ActivityId] = []string{}
		wantDependencies[activities[1].ActivityId] = []string{activities[0].DisplayName}
		wantDependencies[activities[2].ActivityId] = []string{}
		//check if table rows returned is sorted correctly despite adding in reverse order
		trs, err := s.GetTableRowsSortedByRank(t.Context(), p.ProjectId)
		if err != nil {
			t.Fatalf("Failed to get table rows sorted by rank. Encountered error: %v", err)
		}
		for _, tr := range trs {
			if got, want := tr.Dependencies, wantDependencies[tr.ActivityId]; !slices.Equal(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		}
	})

	t.Run("Update Activity Name", func(t *testing.T) {
		s := newStore()
		p := createTestProject(t, s)

		// add activity
		activities := make([]models.Activity, 3)
		activitiesMap := make(map[uuid.UUID]models.Activity)
		for i := range len(activities) {
			activityName := fmt.Sprintf("test activity %v", i+1)
			activityDuration := 5 * time.Minute
			activities[i] = createActivity(activityName, activityDuration, p, t, s)
			activitiesMap[activities[i].ActivityId] = activities[i]
		}

		activityToUpdate := activities[1]
		updatedName := "updated name"

		a, err := s.UpdateActivityName(t.Context(), activityToUpdate.ActivityId, updatedName)
		if err != nil {
			t.Fatalf("Failed to update activity name. Encountered error: %v", err)
		}
		wantUpdatedActivity := models.Activity{ActivityId: activityToUpdate.ActivityId, DisplayName: updatedName, Duration: activityToUpdate.Duration}
		assertActivity(a, wantUpdatedActivity, t)
		wantActivitiesMap := make(map[uuid.UUID]models.Activity)
		maps.Copy(wantActivitiesMap, activitiesMap)
		wantActivitiesMap[wantUpdatedActivity.ActivityId] = wantUpdatedActivity

		getActivities, err := s.GetActivities(t.Context(), p.ProjectId)
		if err != nil {
			t.Fatalf("Failed to get activities. Encountered error %v", err)
		}

		for _, storeActivity := range getActivities {
			assertActivity(storeActivity, wantActivitiesMap[storeActivity.ActivityId], t)
		}
	})

	t.Run("Update Activity Duration", func(t *testing.T) {
		s := newStore()
		p := createTestProject(t, s)

		// add activity
		activities := make([]models.Activity, 3)
		activitiesMap := make(map[uuid.UUID]models.Activity)
		for i := range len(activities) {
			activityName := fmt.Sprintf("test activity %v", i+1)
			activityDuration := 5 * time.Minute
			activities[i] = createActivity(activityName, activityDuration, p, t, s)
			activitiesMap[activities[i].ActivityId] = activities[i]
		}

		activityToUpdate := activities[1]
		updatedDuration := 6 * time.Minute

		a, err := s.UpdateActivityDuration(t.Context(), activityToUpdate.ActivityId, updatedDuration)
		if err != nil {
			t.Fatalf("Failed to update activity duration. Encountered error: %v", err)
		}
		wantUpdatedActivity := models.Activity{ActivityId: activityToUpdate.ActivityId, DisplayName: activityToUpdate.DisplayName, Duration: updatedDuration}
		assertActivity(a, wantUpdatedActivity, t)
		wantActivitiesMap := make(map[uuid.UUID]models.Activity)
		maps.Copy(wantActivitiesMap, activitiesMap)
		wantActivitiesMap[wantUpdatedActivity.ActivityId] = wantUpdatedActivity

		getActivities, err := s.GetActivities(t.Context(), p.ProjectId)
		if err != nil {
			t.Fatalf("Failed to get activities. Encountered error %v", err)
		}

		for _, storeActivity := range getActivities {
			assertActivity(storeActivity, wantActivitiesMap[storeActivity.ActivityId], t)
		}
	})

	t.Run("Update Activity Dependencies", func(t *testing.T) {
		s := newStore()
		p := createTestProject(t, s)

		// add activity
		activities := make([]models.Activity, 3)
		for i := range len(activities) {
			activityName := fmt.Sprintf("test activity %v", i+1)
			activityDuration := 5 * time.Minute
			activities[i] = createActivity(activityName, activityDuration, p, t, s)
		}
		//add table rows
		for i := range len(activities) {
			activityToAdd := activities[i]
			createTableRow(p, activityToAdd, int64(i), t, s)
		}

		a, err := s.UpdateActivityDependencies(t.Context(), p.ProjectId, activities[1].ActivityId, []uuid.UUID{activities[0].ActivityId})
		if err != nil {
			t.Fatalf("Failed to create dependency. Encountered error: %v", err)
		}
		if got, want := a[0].ActivityId, activities[0].ActivityId; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}

		wantDependencies := map[uuid.UUID][]string{}
		wantDependencies[activities[0].ActivityId] = []string{}
		wantDependencies[activities[1].ActivityId] = []string{activities[0].DisplayName}
		wantDependencies[activities[2].ActivityId] = []string{}
		//check if table rows returned is sorted correctly despite adding in reverse order
		trs, err := s.GetTableRowsSortedByRank(t.Context(), p.ProjectId)
		if err != nil {
			t.Fatalf("Failed to get table rows sorted by rank. Encountered error: %v", err)
		}
		for _, tr := range trs {
			if got, want := tr.Dependencies, wantDependencies[tr.ActivityId]; !slices.Equal(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		}
	})
}

func TestMockStoreContract(t *testing.T) {
	testModelContract(t, func() models.Model {
		return NewMockStore()
	})
}

func TestSqliteStoreContract(t *testing.T) {
	testModelContract(t, func() models.Model {
		return NewSqliteStore(sqlitedb.New(sqlitedb.NewSQLiteDb(":memory:")))
	})
}
