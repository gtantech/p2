package db

import (
	"context"
	"slices"
	"testing"
)

type mock_dependency struct {
	ID                    string
	ProjectID             string
	Relationship          string
	PredecessorActivityID string
	SuccessorActivityID   string
}

type mock_activity struct {
	ID        string
	ProjectID string
	DispName  string
	Duration  int64
}

func TestFindAllPredecessorNamesBySuccessor(t *testing.T) {
	database, err := NewSQLiteStorage(":memory:")
	if err != nil {
		t.Errorf("failed to open database with error: %v", err)
	}
	activity1 := mock_activity{
		ID:        "test_activity1_id",
		ProjectID: "test_project1_id",
		DispName:  "activity1",
		Duration:  0,
	}

	activity2 := mock_activity{
		ID:        "test_activity2_id",
		ProjectID: activity1.ProjectID,
		DispName:  "activity2",
		Duration:  0,
	}

	activity3 := mock_activity{
		ID:        "test_activity3_id",
		ProjectID: activity1.ProjectID,
		DispName:  "activity3",
		Duration:  0,
	}

	activities := []mock_activity{activity1, activity2, activity3}

	for _, activity := range activities {
		_, err := database.Exec("INSERT INTO activities (id, project_id, disp_name, duration) VALUES (?, ?, ?, ?)",
			activity.ID, activity.ProjectID, activity.DispName, activity.Duration)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}

	// a1 -> a2
	// a3 -> a2

	dependency1 := mock_dependency{
		ID:                    "test_id1",
		ProjectID:             activity1.ProjectID,
		Relationship:          "SS",
		PredecessorActivityID: activity1.ID,
		SuccessorActivityID:   activity2.ID,
	}
	dependency2 := mock_dependency{
		ID:                    "test_id2",
		ProjectID:             activity1.ProjectID,
		Relationship:          "FS",
		PredecessorActivityID: activity3.ID,
		SuccessorActivityID:   activity2.ID,
	}

	dependencies := []mock_dependency{dependency1, dependency2}

	for _, dependency := range dependencies {
		_, err := database.Exec("INSERT INTO dependencies (id, project_id, relationship, predecessor_activity_id, successor_activity_id) VALUES (?, ?, ?, ?, ?)",
			dependency.ID, dependency.ProjectID, dependency.Relationship, dependency.PredecessorActivityID, dependency.SuccessorActivityID)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}

	d := New(database)

	results, err := d.FindAllPredecessorNamesBySuccessor(context.Background(), activity2.ID)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if got, want, want2 := []string{results[0].PredecessorActivityID, results[1].PredecessorActivityID}, []string{activity1.ID, activity3.ID}, []string{activity3.ID, activity1.ID}; !(slices.Equal(got, want) || slices.Equal(got, want2)) {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want, want2 := []string{results[0].PredecessorActivityName, results[1].PredecessorActivityName}, []string{activity1.DispName, activity3.DispName}, []string{activity3.DispName, activity1.DispName}; !(slices.Equal(got, want) || slices.Equal(got, want2)) {
		t.Errorf("got %v, want %v", got, want)
	}
}
