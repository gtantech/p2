package store

import (
	"testing"
	"uuid"

	"github.com/gtantech/p2/internal/db"
)

func TestNewDependency(t *testing.T) {
	dependencyId := uuid.New()
	projectId := uuid.New()
	relationship := SS
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
	relationship := SS
	predecessorId := uuid.New()
	successorId := uuid.New()

	params := CreateDepdencencyParams{
		ProjectID:             projectId,
		Relationship:          relationship,
		PredecessorActivityID: predecessorId,
		SuccessorActivityID:   successorId,
	}

	dbParams := params.ToDbInsertDependencyParams(dependencyId)

	if got, want := uuid.MustParse(dbParams.ID), dependencyId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := uuid.MustParse(dbParams.ProjectID), projectId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := RelationshipType(dbParams.Relationship), relationship; got != want {
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
	relationship := SS
	predecessorId := uuid.New()
	successorId := uuid.New()

	params := UpdateDepdencencyParams{
		ID:                    dependencyId,
		Relationship:          relationship,
		PredecessorActivityID: predecessorId,
		SuccessorActivityID:   successorId,
	}

	dbParams := params.ToDbUpdateDependencyParams()

	if got, want := uuid.MustParse(dbParams.ID), dependencyId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := RelationshipType(dbParams.Relationship), relationship; got != want {
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
