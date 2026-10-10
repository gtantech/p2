package store

import (
	"cmp"
	"slices"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models"
)

type MockStore struct {
	activities       []models.Activity
	activitiesLookup map[uuid.UUID]*models.Activity
	dependencies     map[uuid.UUID][]*models.Activity
	tableRows        []*models.TableRow
	tableRowsLookup  map[uuid.UUID]*models.TableRow
}

// GetNextTableRowByActivityId implements [models.Model].
func (ms *MockStore) GetNextTableRowByActivityId(projectId uuid.UUID, activityId uuid.UUID) (models.TableRow, error) {
	if _, ok := ms.activitiesLookup[activityId]; !ok {
		return models.TableRow{}, ErrActivityNotFound
	}

	trs, err := ms.GetTableRowsSortedByRank(projectId)
	if err != nil {
		return models.TableRow{}, err
	}
	for i := range trs {
		if i == 0 {
			continue
		}
		if trs[i-1].ActivityId == activityId {
			return trs[i], nil
		}
	}
	return models.TableRow{}, ErrEndOfTable
}

// GetTableRowByActivityId implements [models.Model].
func (ms *MockStore) GetTableRowByActivityId(activityId uuid.UUID) (models.TableRow, error) {
	if _, ok := ms.activitiesLookup[activityId]; !ok {
		return models.TableRow{}, ErrActivityNotFound
	}
	return *ms.tableRowsLookup[activityId], nil
}

// UpdateActivityDependencies implements [models.Model].
func (ms *MockStore) UpdateActivityDependencies(activityId uuid.UUID, activityDependencyNames []string) ([]models.Activity, error) {
	if _, ok := ms.activitiesLookup[activityId]; !ok {
		return []models.Activity{}, ErrActivityNotFound
	}
	dependenciesAddr := []*models.Activity{}
	dependencies := []models.Activity{}

	for _, d := range activityDependencyNames {
		//find in activities
		for j, a := range ms.activities {
			if a.DisplayName == d {
				dependencies = append(dependencies, ms.activities[j])
				dependenciesAddr = append(dependenciesAddr, &ms.activities[j])
				break
			}
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
func (ms *MockStore) UpdateActivityDuration(activityId uuid.UUID, duration time.Duration) (models.Activity, error) {
	if _, ok := ms.activitiesLookup[activityId]; !ok {
		return models.Activity{}, ErrActivityNotFound
	}
	ms.activitiesLookup[activityId].Duration = duration
	return *ms.activitiesLookup[activityId], nil
}

// UpdateActivityName implements [models.Model].
func (ms *MockStore) UpdateActivityName(activityId uuid.UUID, activityName string) (models.Activity, error) {
	if _, ok := ms.activitiesLookup[activityId]; !ok {
		return models.Activity{}, ErrActivityNotFound
	}
	ms.activitiesLookup[activityId].DisplayName = activityName
	return *ms.activitiesLookup[activityId], nil
}

// GetTableRowsSortedByRank implements [models.Model].
func (ms *MockStore) GetTableRowsSortedByRank(projectId uuid.UUID) ([]models.TableRow, error) {
	trs := ms.tableRows
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
		}
		return trsSorted
	}(), nil
}

// CreateTableRow implements [models.Model].
func (ms *MockStore) CreateTableRow(projectId uuid.UUID, activityId uuid.UUID, rowSortRank int64) (models.TableRow, error) {
	if _, ok := ms.activitiesLookup[activityId]; !ok {
		return models.TableRow{}, ErrActivityNotFound
	}
	dependenciesStr := []string{}
	dependencies := ms.dependencies[activityId]
	for _, d := range dependencies {
		dependenciesStr = append(dependenciesStr, d.DisplayName)
	}
	tr := models.TableRow{Activity: ms.activitiesLookup[activityId], ProjectId: projectId, Dependencies: dependenciesStr, SortRank: rowSortRank}
	ms.tableRows = append(ms.tableRows, &tr)
	ms.tableRowsLookup[activityId] = ms.tableRows[len(ms.tableRows)-1]
	return tr, nil
}

// CreateActivity implements [models.Model].
func (ms *MockStore) CreateActivity(projectId uuid.UUID, activityName string, duration time.Duration) (models.Activity, error) {
	a := models.Activity{ActivityId: uuid.NewV7(), DisplayName: activityName, Duration: duration}
	ms.activitiesLookup[a.ActivityId] = &a
	ms.activities = append(ms.activities, a)
	return a, nil
}

func NewMockStore() *MockStore {
	ms := &MockStore{}
	ms.activitiesLookup = map[uuid.UUID]*models.Activity{}
	ms.tableRowsLookup = map[uuid.UUID]*models.TableRow{}
	ms.activities = []models.Activity{
		{ActivityId: uuid.MustParse("dcce6c98-b31b-4ed6-a2ed-15d24ae96b41"), DisplayName: "A", Duration: 5 * time.Minute},
		{ActivityId: uuid.MustParse("a4518b28-e597-4295-a532-3a501a75ab2b"), DisplayName: "B", Duration: 4 * time.Minute},
		{ActivityId: uuid.MustParse("aa5a833c-2407-464c-93c4-fa9edba03b52"), DisplayName: "C", Duration: 3 * time.Minute},
	}

	ms.dependencies = map[uuid.UUID][]*models.Activity{}
	ms.dependencies[ms.activities[1].ActivityId] = []*models.Activity{&ms.activities[0]}
	ms.dependencies[ms.activities[2].ActivityId] = []*models.Activity{&ms.activities[0], &ms.activities[1]}

	ms.tableRows = []*models.TableRow{}

	for i, a := range ms.activities {
		ms.activitiesLookup[a.ActivityId] = &ms.activities[i]
		ms.CreateTableRow(uuid.Max(), a.ActivityId, int64(i)*models.TableRowSortRankStep)
	}
	return ms
}

// GetActivities implements [models.Model].
func (ms *MockStore) GetActivities(projectId uuid.UUID) ([]models.Activity, error) {
	return ms.activities, nil
}

// GetHome implements [models.Model].
func (ms *MockStore) GetHome() models.Home {
	return models.Home{Name: "World"}
}

var _ models.Model = (*MockStore)(nil) //ensures MockStore implements models.Model at compile time
