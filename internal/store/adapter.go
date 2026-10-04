package store

import (
	"context"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models/viewmodels"
	"github.com/gtantech/p2/internal/routes"
)

type StoreViewMockAdapter struct {
	table map[uuid.UUID]viewmodels.Table
}

// CreateDependencyTableRow implements [routes.StoreView].
func (s *StoreViewMockAdapter) CreateDependencyTableRow(projectId uuid.UUID, activityName string, sortRank int64, ctx context.Context) (viewmodels.TableRow, error) {
	return viewmodels.TableRow{ActivityId: uuid.New(), ActivityName: activityName, ProjectId: projectId, SortRank: sortRank}, nil
}

// CreateEmptyDependencyTableRow implements [routes.StoreView].
func (s *StoreViewMockAdapter) CreateEmptyDependencyTableRow(projectId uuid.UUID, sortRank int64, ctx context.Context) (viewmodels.TableRow, error) {
	return viewmodels.TableRow{ActivityId: uuid.New(), ProjectId: projectId, SortRank: sortRank}, nil
}

func NewStoreViewMockAdapter() *StoreViewMockAdapter {
	s := StoreViewMockAdapter{
		table: make(map[uuid.UUID]viewmodels.Table),
	}
	projectId := uuid.Max()
	s.table[projectId] = viewmodels.Table{ProjectId: projectId, Rows: []viewmodels.TableRow{viewmodels.TableRow{ActivityId: uuid.New(), ProjectId: projectId, ActivityName: "A", PredecessorActivities: []string{"B", "C"}, Duration: 5 * time.Minute, SortRank: 0}}}

	return &s
}

// GetDependencyTableByProjectId implements [routes.StoreView].
func (s *StoreViewMockAdapter) GetDependencyTableByProjectId(projectId uuid.UUID, ctx context.Context) (viewmodels.Table, error) {
	return s.table[projectId], nil
}

var _ routes.StoreView = (*StoreViewMockAdapter)(nil) //ensures ExampleStruct implements ExampleInterface at compile time
