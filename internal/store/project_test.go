package store

import (
	"testing"
	"uuid"

	"github.com/gtantech/p2/internal/db"
)

func TestNewProject(t *testing.T) {
	projectId := uuid.New()
	displayName := "test_disp_name"

	p := newProject(projectId, displayName)

	if got, want := p.ID, projectId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := p.DisplayName, displayName; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestToDbInsertProjectParams(t *testing.T) {
	projectId := uuid.New()
	displayName := "test_disp_name"

	params := CreateProjectParams{
		DisplayName: displayName,
	}

	dbParams := params.toDbInsertProjectParams(projectId)

	if got, want := uuid.MustParse(dbParams.ID), projectId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := dbParams.DispName, displayName; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestToDbUpdateProjectParams(t *testing.T) {
	projectId := uuid.New()
	displayName := "test_disp_name"

	params := UpdateProjectParams{
		Id:          projectId,
		DisplayName: displayName,
	}

	dbParams := params.toDbUpdateProjectParams()

	if got, want := uuid.MustParse(dbParams.ID), projectId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := dbParams.DispName, displayName; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNewProjectStoreFromDb(t *testing.T) {
	var q = &db.Queries{}
	s := NewProjectStoreFromDb(q)

	if got, want := s.queries, q; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
