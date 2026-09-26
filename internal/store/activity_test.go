package store

import (
	"testing"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/db"
	"github.com/gtantech/p2/internal/models"
)

func TestNewActivity(t *testing.T) {
	activityId := uuid.NewV7()
	projectId := uuid.NewV7()
	displayName := "test_disp_name"
	duration := 5 * time.Minute
	a := newActivity(activityId, projectId, displayName, duration)

	if got, want := a.ID, activityId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := a.ProjectID, projectId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := a.DisplayName, displayName; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := time.Duration(a.Duration), duration; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

}

func TestToDbInsertActivityParams(t *testing.T) {
	activityId := uuid.NewV7()
	projectId := uuid.NewV7()
	displayName := "test_disp_name"
	duration := 5 * time.Minute

	params := models.CreateActivityParams{ProjectID: projectId, DisplayName: displayName, Duration: duration}
	dbParams := toDbInsertActivityParams(activityId, &params)

	if got, want := uuid.MustParse(dbParams.ID), activityId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := uuid.MustParse(dbParams.ProjectID), projectId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := dbParams.DispName, displayName; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := time.Duration(dbParams.Duration), duration; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestToDbFindAllActivitiesByNameAndProjectParams(t *testing.T) {
	projectId := uuid.NewV7()
	displayName := "test_disp_name"

	params := models.GetActivityByNameAndProjectParams{ProjectID: projectId, DisplayName: displayName}
	dbParams := toDbFindAllActivitiesByNameAndProjectParams(&params)

	if got, want := dbParams.DispName, displayName; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := uuid.MustParse(dbParams.ProjectID), projectId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestToUpdateActivityParams(t *testing.T) {
	activityId := uuid.NewV7()
	displayName := "test_disp_name"
	duration := 5 * time.Minute

	params := models.UpdateActivityParams{Id: activityId, DisplayName: displayName, Duration: duration}
	dbParams := toUpdateActivityParams(&params)

	if got, want := uuid.MustParse(dbParams.ID), activityId; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := dbParams.DispName, displayName; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := time.Duration(dbParams.Duration), duration; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNewActivityStoreFromDb(t *testing.T) {
	var q = &db.Queries{}
	s := NewActivityStoreFromDb(q)

	if got, want := s.queries, q; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
