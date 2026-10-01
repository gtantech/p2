package sqlitedb

import (
	"context"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models/viewmodels"
	"github.com/gtantech/p2/internal/routes"
)

type StoreViewSqliteAdapter struct {
	queries *Queries
}

// CreateEmptyDependencyTableRow implements [routes.StoreView].
func (s *StoreViewSqliteAdapter) CreateEmptyDependencyTableRow(projectId uuid.UUID, sortRank int64, ctx context.Context) (viewmodels.TableRow, error) {
	activity, err := s.queries.InsertActivity(ctx, InsertActivityParams{
		ID:        uuid.NewV7().String(),
		ProjectID: projectId.String(),
		DispName:  "",
		Duration:  0,
	})
	s.queries.InsertActivityOrdering(ctx, InsertActivityOrderingParams{uuid.NewV7().String(), activity.ProjectID, activity.ID, sortRank})
	if err != nil {
		return viewmodels.TableRow{}, err
	}
	return viewmodels.TableRow{ActivityId: uuid.MustParse(activity.ID), ProjectId: projectId, ActivityName: activity.DispName, PredecessorActivities: []string{}, Duration: time.Duration(activity.Duration), SortRank: sortRank}, nil
}

// GetDependencyTableByProjectId implements [routes.StoreView].
func (s *StoreViewSqliteAdapter) GetDependencyTableByProjectId(projectId uuid.UUID, ctx context.Context) (viewmodels.Table, error) {
	successorActivities, err := s.queries.FindAllActivitiesByProjectSorted(ctx, projectId.String())
	if err != nil {
		return viewmodels.Table{}, err
	}
	dependencies, err := s.queries.FindAllDependenciesByProject(ctx, projectId.String())
	if err != nil {
		return viewmodels.Table{}, err
	}
	if len(successorActivities) == 0 {
		//return early with a table with no rows
		return viewmodels.Table{ProjectId: projectId, Rows: make([]viewmodels.TableRow, 0)}, nil
	}

	successorActivityMap := map[string]FindAllActivitiesByProjectSortedRow{}
	tableRows := make([]viewmodels.TableRow, len(successorActivities))

	//add viewmodels.Activities to table rows
	//add successor activites to successorActivityMap where the successor id is the key and the successor activity is the value
	for i, successorActivity := range successorActivities {
		successorActivityMap[successorActivity.ID] = successorActivity
		tableRows[i] = viewmodels.TableRow{
			ActivityId:            uuid.MustParse(successorActivity.ID),
			ProjectId:             projectId,
			ActivityName:          successorActivity.DispName,
			PredecessorActivities: make([]string, 0),
			Duration:              time.Duration(successorActivity.Duration),
			SortRank:              successorActivity.SortRank,
		}
	}

	//create a dependency map where the successor id is the key and the dependencies are the values
	dependencyMap := map[string][]string{}
	for _, dependency := range dependencies {
		//get the predecessor based on the activity map with the predecessor id as key
		predecessor := successorActivityMap[dependency.PredecessorActivityID]
		//append predecessor name to dependency map with successor.id as key
		dependencyMap[dependency.SuccessorActivityID] = append(dependencyMap[dependency.SuccessorActivityID], predecessor.DispName)
	}

	//update PredecessorActivities field
	for i := range tableRows {
		tableRows[i].PredecessorActivities = dependencyMap[tableRows[i].ActivityId.String()]
	}
	return viewmodels.Table{ProjectId: projectId, Rows: tableRows}, nil
}

func NewStoreViewSqliteAdapter(queries *Queries) *StoreViewSqliteAdapter {
	s := StoreViewSqliteAdapter{
		queries: queries,
	}

	return &s
}

var _ routes.StoreView = (*StoreViewSqliteAdapter)(nil) //ensures ExampleStruct implements ExampleInterface at compile time

type StoreFormSqliteAdapter struct {
	queries *Queries
}

// UpdateActivityDuration implements [routes.Store].
func (s *StoreFormSqliteAdapter) UpdateActivityDuration(activityId uuid.UUID, activityDuration time.Duration, ctx context.Context) error {
	_, err := s.queries.UpdateActivityDuration(ctx, UpdateActivityDurationParams{Duration: int64(activityDuration), ID: activityId.String()})
	return err
}

// UpdateActivityName implements [routes.Store].
func (s *StoreFormSqliteAdapter) UpdateActivityName(activityId uuid.UUID, activityName string, ctx context.Context) error {
	_, err := s.queries.UpdateActivityName(ctx, UpdateActivityNameParams{DispName: activityName, ID: activityId.String()})
	return err
}

func NewStoreFormSqliteAdapter(queries *Queries) *StoreFormSqliteAdapter {
	return &StoreFormSqliteAdapter{queries: queries}
}

var _ routes.Store = (*StoreFormSqliteAdapter)(nil) //ensures ExampleStruct implements ExampleInterface at compile time
