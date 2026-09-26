package store

import (
	"slices"
	"testing"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/db"
	"github.com/gtantech/p2/internal/models"
)

type mock_dependency struct {
	ID                    uuid.UUID
	ProjectID             uuid.UUID
	Relationship          string
	PredecessorActivityID uuid.UUID
	SuccessorActivityID   uuid.UUID
}

type mock_activity struct {
	ID        uuid.UUID
	ProjectID uuid.UUID
	DispName  string
	Duration  time.Duration
}

func TestNewDependency(t *testing.T) {
	dependencyId := uuid.New()
	projectId := uuid.New()
	relationship := models.SS
	predecessorId := uuid.New()
	successorId := uuid.New()

	d := newDependency(dependencyId, projectId, relationship, predecessorId, successorId)

	if got, want := d.ID, dependencyId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := d.ProjectID, projectId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := d.Relationship, relationship; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := d.PredecessorActivityID, predecessorId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := d.SuccessorActivityID, successorId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestToDbInsertDependencyParams(t *testing.T) {
	dependencyId := uuid.New()
	projectId := uuid.New()
	relationship := models.SS
	predecessorId := uuid.New()
	successorId := uuid.New()

	params := models.CreateDepdencencyParams{
		ProjectID:             projectId,
		Relationship:          relationship,
		PredecessorActivityID: predecessorId,
		SuccessorActivityID:   successorId,
	}

	dbParams := toDbInsertDependencyParams(dependencyId, &params)

	if got, want := uuid.MustParse(dbParams.ID), dependencyId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := uuid.MustParse(dbParams.ProjectID), projectId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := models.RelationshipType(dbParams.Relationship), relationship; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := uuid.MustParse(dbParams.PredecessorActivityID), predecessorId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := uuid.MustParse(dbParams.SuccessorActivityID), successorId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestToDbUpdateDependencyParams(t *testing.T) {
	dependencyId := uuid.New()
	relationship := models.SS
	predecessorId := uuid.New()
	successorId := uuid.New()

	params := models.UpdateDepdencencyParams{
		ID:                    dependencyId,
		Relationship:          relationship,
		PredecessorActivityID: predecessorId,
		SuccessorActivityID:   successorId,
	}

	dbParams := toDbUpdateDependencyParams(&params)

	if got, want := uuid.MustParse(dbParams.ID), dependencyId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := models.RelationshipType(dbParams.Relationship), relationship; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := uuid.MustParse(dbParams.PredecessorActivityID), predecessorId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := uuid.MustParse(dbParams.SuccessorActivityID), successorId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNewDependencyStoreFromDb(t *testing.T) {
	var q = &db.Queries{}
	s := NewDependencyStoreFromDb(q)

	if got, want := s.queries, q; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetPredecessorNamesBySuccessor(t *testing.T) {
	database, err := db.NewSQLiteStorage(":memory:")
	if err != nil {
		t.Errorf("failed to open database with error: %v", err)
	}
	projectId := uuid.New()
	activity1 := mock_activity{
		ID:        uuid.New(),
		ProjectID: projectId,
		DispName:  "activity1",
		Duration:  0,
	}

	activity2 := mock_activity{
		ID:        uuid.New(),
		ProjectID: activity1.ProjectID,
		DispName:  "activity2",
		Duration:  0,
	}

	activity3 := mock_activity{
		ID:        uuid.New(),
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
		ID:                    uuid.New(),
		ProjectID:             activity1.ProjectID,
		Relationship:          "SS",
		PredecessorActivityID: activity1.ID,
		SuccessorActivityID:   activity2.ID,
	}
	dependency2 := mock_dependency{
		ID:                    uuid.New(),
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

	store := NewStoreFromDb(db.New(database))

	results, err := store.Dependency.GetPredecessorNamesBySuccessor(t.Context(), activity2.ID)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if got, want, want2 := []uuid.UUID{results[0].PredecessorActivityID, results[1].PredecessorActivityID}, []uuid.UUID{activity1.ID, activity3.ID}, []uuid.UUID{activity3.ID, activity1.ID}; !(slices.Equal(got, want) || slices.Equal(got, want2)) {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want, want2 := []string{results[0].PredecessorActivityName, results[1].PredecessorActivityName}, []string{activity1.DispName, activity3.DispName}, []string{activity3.DispName, activity1.DispName}; !(slices.Equal(got, want) || slices.Equal(got, want2)) {
		t.Errorf("got %v, want %v", got, want)
	}
}
