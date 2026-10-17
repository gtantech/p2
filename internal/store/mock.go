package store

import (
	"cmp"
	"context"
	"errors"
	"slices"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models"
)

var _ models.Model = (*MockStore)(nil) //ensures MockStore implements models.Model at compile time

type project struct {
	projectName string
	activities  []*models.Activity
	tableRows   []*models.TableRow
}

type MockStore struct {
	projectLookup    map[uuid.UUID]*project
	activitiesLookup map[uuid.UUID]*models.Activity
	tableRowsLookup  map[uuid.UUID]*models.TableRow
	dependencies     map[uuid.UUID][]*models.Activity
}

// CreateProject implements [models.Model].
func (ms *MockStore) CreateProject(ctx context.Context, projectName string) (models.Project, error) {
	p := models.Project{
		ProjectId:          uuid.NewV7(),
		ProjectDisplayName: projectName,
	}
	ms.projectLookup[p.ProjectId] = &project{
		projectName: p.ProjectDisplayName,
		activities:  []*models.Activity{},
		tableRows:   []*models.TableRow{},
	}
	return p, nil
}

// CreateDependency implements [models.Model].
func (ms *MockStore) CreateDependency(ctx context.Context, projectId uuid.UUID, successorActivityId uuid.UUID, predecessorActivityId uuid.UUID) (models.Dependency, error) {
	if _, ok := ms.projectLookup[projectId]; !ok {
		return models.Dependency{}, ErrProjectNotFound
	}
	if _, ok := ms.activitiesLookup[successorActivityId]; !ok {
		return models.Dependency{}, ErrActivityNotFound
	}
	if _, ok := ms.activitiesLookup[predecessorActivityId]; !ok {
		return models.Dependency{}, ErrActivityNotFound
	}
	ms.dependencies[successorActivityId] = append(ms.dependencies[successorActivityId], ms.activitiesLookup[predecessorActivityId])
	return models.Dependency{DependencyId: successorActivityId, SuccessorActivityId: successorActivityId, PredecessorActivityId: predecessorActivityId}, nil
}

// CreateTableRowBelow implements [models.Model].
func (ms *MockStore) CreateTableRowBelow(ctx context.Context, projectId uuid.UUID, activityId uuid.UUID, belowActivityId uuid.UUID) (models.TableRow, error) {
	if _, ok := ms.projectLookup[projectId]; !ok {
		return models.TableRow{}, ErrProjectNotFound
	}
	if _, ok := ms.activitiesLookup[activityId]; !ok {
		return models.TableRow{}, ErrActivityNotFound
	}
	if _, ok := ms.activitiesLookup[belowActivityId]; !ok {
		return models.TableRow{}, ErrActivityNotFound
	}

	rowAfter, err := func() (models.TableRow, error) {
		trs, err := ms.GetTableRowsSortedByRank(ctx, projectId)
		if err != nil {
			return models.TableRow{}, err
		}
		for i := range trs {
			if i == 0 {
				continue
			}
			if trs[i-1].ActivityId == belowActivityId {
				return trs[i], nil
			}
		}
		return models.TableRow{}, ErrEndOfTable
	}()

	insertRowRank := int64(0)
	rowCurrent := ms.tableRowsLookup[belowActivityId]
	if err != nil {
		if errors.Is(err, ErrEndOfTable) {
			insertRowRank = rowCurrent.SortRank + models.TableRowSortRankStep
		} else {
			return models.TableRow{}, err
		}
	} else {
		insertRowRank = (rowCurrent.SortRank / 2) + (rowAfter.SortRank / 2)
	}

	return ms.CreateTableRow(ctx, projectId, activityId, insertRowRank)
}

// CreateTableRowAtBottom implements [models.Model].
func (ms *MockStore) CreateTableRowAtBottom(ctx context.Context, projectId uuid.UUID, activityId uuid.UUID) (models.TableRow, error) {
	if _, ok := ms.projectLookup[projectId]; !ok {
		return models.TableRow{}, ErrProjectNotFound
	}
	maxSortRank := int64(0)
	for _, tr := range ms.projectLookup[projectId].tableRows {
		if tr.SortRank > maxSortRank {
			maxSortRank = tr.SortRank
		}
	}

	return ms.CreateTableRow(ctx, projectId, activityId, maxSortRank+models.TableRowSortRankStep)
}

// UpdateActivityDependencies implements [models.Model].
func (ms *MockStore) UpdateActivityDependencies(ctx context.Context, projectId uuid.UUID, activityId uuid.UUID, dependencyActivityIds []uuid.UUID) ([]models.Activity, error) {
	if _, ok := ms.activitiesLookup[activityId]; !ok {
		return []models.Activity{}, ErrActivityNotFound
	}
	dependenciesAddr := []*models.Activity{}
	dependencies := []models.Activity{}

	for _, d := range dependencyActivityIds {
		if a, ok := ms.activitiesLookup[d]; ok {
			dependencies = append(dependencies, *a)
			dependenciesAddr = append(dependenciesAddr, a)
		} else {
			return []models.Activity{}, ErrActivityNotFound
		}
	}

	ms.tableRowsLookup[activityId].Dependencies = func() []string {
		dependencyStrSlice := make([]string, len(dependencies))
		for i := range dependencies {
			dependencyStrSlice[i] = dependencies[i].DisplayName
		}
		return dependencyStrSlice
	}()
	ms.dependencies[activityId] = dependenciesAddr
	return dependencies, nil
}

// UpdateActivityDuration implements [models.Model].
func (ms *MockStore) UpdateActivityDuration(ctx context.Context, activityId uuid.UUID, duration time.Duration) (models.Activity, error) {
	if _, ok := ms.activitiesLookup[activityId]; !ok {
		return models.Activity{}, ErrActivityNotFound
	}
	ms.activitiesLookup[activityId].Duration = duration
	return *ms.activitiesLookup[activityId], nil
}

// UpdateActivityName implements [models.Model].
func (ms *MockStore) UpdateActivityName(ctx context.Context, activityId uuid.UUID, activityName string) (models.Activity, error) {
	if _, ok := ms.activitiesLookup[activityId]; !ok {
		return models.Activity{}, ErrActivityNotFound
	}
	ms.activitiesLookup[activityId].DisplayName = activityName
	return *ms.activitiesLookup[activityId], nil
}

// GetTableRowsSortedByRank implements [models.Model].
func (ms *MockStore) GetTableRowsSortedByRank(ctx context.Context, projectId uuid.UUID) ([]models.TableRow, error) {
	if _, ok := ms.projectLookup[projectId]; !ok {
		return []models.TableRow{}, ErrProjectNotFound
	}
	trs := ms.projectLookup[projectId].tableRows
	slices.SortFunc(trs, func(a, b *models.TableRow) int {
		if c := cmp.Compare(a.SortRank, b.SortRank); c != 0 {
			return c
		}
		return cmp.Compare(a.SortRank, b.SortRank)
	})
	return func() []models.TableRow {
		trsSorted := make([]models.TableRow, len(trs))
		for i := range trs {
			trsSorted[i] = *trs[i]
			trsSorted[i].Dependencies = func() []string {
				dependenciesStr := []string{}
				dependencies := ms.dependencies[trs[i].ActivityId]
				for _, d := range dependencies {
					dependenciesStr = append(dependenciesStr, d.DisplayName)
				}
				return dependenciesStr
			}()
		}
		return trsSorted
	}(), nil
}

// CreateTableRow implements [models.Model].
func (ms *MockStore) CreateTableRow(ctx context.Context, projectId uuid.UUID, activityId uuid.UUID, rowSortRank int64) (models.TableRow, error) {
	if _, ok := ms.projectLookup[projectId]; !ok {
		return models.TableRow{}, ErrProjectNotFound
	}
	if _, ok := ms.activitiesLookup[activityId]; !ok {
		return models.TableRow{}, ErrActivityNotFound
	}
	dependenciesStr := []string{}
	dependencies := ms.dependencies[activityId]
	for _, d := range dependencies {
		dependenciesStr = append(dependenciesStr, d.DisplayName)
	}
	tr := models.TableRow{Activity: ms.activitiesLookup[activityId], ProjectId: projectId, Dependencies: dependenciesStr, SortRank: rowSortRank}
	ms.projectLookup[projectId].tableRows = append(ms.projectLookup[projectId].tableRows, &tr)
	ms.tableRowsLookup[activityId] = ms.projectLookup[projectId].tableRows[len(ms.projectLookup[projectId].tableRows)-1]
	return tr, nil
}

// CreateActivity implements [models.Model].
func (ms *MockStore) CreateActivity(ctx context.Context, projectId uuid.UUID, activityName string, duration time.Duration) (models.Activity, error) {
	if _, ok := ms.projectLookup[projectId]; !ok {
		return models.Activity{}, ErrProjectNotFound
	}
	a := models.Activity{ActivityId: uuid.NewV7(), DisplayName: activityName, Duration: duration}
	ms.activitiesLookup[a.ActivityId] = &a
	ms.projectLookup[projectId].activities = append(ms.projectLookup[projectId].activities, &a)
	return a, nil
}

func NewMockStore() *MockStore {
	ms := &MockStore{}
	ms.projectLookup = make(map[uuid.UUID]*project)
	ms.activitiesLookup = map[uuid.UUID]*models.Activity{}
	ms.tableRowsLookup = map[uuid.UUID]*models.TableRow{}
	ms.dependencies = map[uuid.UUID][]*models.Activity{}
	return ms
}

func NewMockStoreWithValues() *MockStore {
	ms := NewMockStore()
	ms.projectLookup[uuid.Max()] = &project{
		projectName: "mock project",
		activities:  []*models.Activity{},
		tableRows:   []*models.TableRow{},
	}
	ms.projectLookup[uuid.Max()].activities = []*models.Activity{
		{ActivityId: uuid.MustParse("dcce6c98-b31b-4ed6-a2ed-15d24ae96b41"), DisplayName: "A", Duration: 5 * time.Minute},
		{ActivityId: uuid.MustParse("a4518b28-e597-4295-a532-3a501a75ab2b"), DisplayName: "B", Duration: 4 * time.Minute},
		{ActivityId: uuid.MustParse("aa5a833c-2407-464c-93c4-fa9edba03b52"), DisplayName: "C", Duration: 3 * time.Minute},
	}

	ms.dependencies[ms.projectLookup[uuid.Max()].activities[1].ActivityId] = []*models.Activity{ms.projectLookup[uuid.Max()].activities[0]}
	ms.dependencies[ms.projectLookup[uuid.Max()].activities[2].ActivityId] = []*models.Activity{ms.projectLookup[uuid.Max()].activities[0], ms.projectLookup[uuid.Max()].activities[1]}

	ms.projectLookup[uuid.Max()].tableRows = []*models.TableRow{}

	for i, a := range ms.projectLookup[uuid.Max()].activities {
		ms.activitiesLookup[a.ActivityId] = ms.projectLookup[uuid.Max()].activities[i]
		ms.CreateTableRow(context.Background(), uuid.Max(), a.ActivityId, int64(i)*models.TableRowSortRankStep)
	}
	return ms
}

// GetActivities implements [models.Model].
func (ms *MockStore) GetActivities(ctx context.Context, projectId uuid.UUID) ([]models.Activity, error) {
	if _, ok := ms.projectLookup[projectId]; !ok {
		return []models.Activity{}, ErrProjectNotFound
	}
	return func() []models.Activity {
		activities := make([]models.Activity, len(ms.projectLookup[projectId].activities))
		for i, a := range ms.projectLookup[projectId].activities {
			activities[i] = *a
		}
		return activities
	}(), nil
}
