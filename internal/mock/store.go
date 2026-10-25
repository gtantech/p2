package mock

import (
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/store"
)

var _ models.Store = (*Store)(nil) //ensures MockStore implements models.Store at compile time

func NewStore() *Store {
	s := &Store{project: make(map[uuid.UUID]home)}
	mockProjectId := uuid.Max()
	mockTable := models.Table{
		Rows: []models.TableRow{
			{RowId: uuid.MustParse("dcce6c98-b31b-4ed6-a2ed-15d24ae96b41"), ActivityName: "A", DependencyNames: []string{}, Duration: 5 * time.Minute},
			{RowId: uuid.MustParse("a4518b28-e597-4295-a532-3a501a75ab2b"), ActivityName: "B", DependencyNames: []string{}, Duration: 2 * time.Minute},
			{RowId: uuid.MustParse("aa5a833c-2407-464c-93c4-fa9edba03b52"), ActivityName: "C", DependencyNames: []string{"A", "B"}, Duration: 3 * time.Minute},
		},
	}
	s.project[mockProjectId] = home{Project: models.Project{ProjecId: mockProjectId, ProjectName: "mock project"}, Table: mockTable}
	return s
}

type Store struct {
	project map[uuid.UUID]home
}

// GetHomeProject implements [models.Store].
func (ms *Store) GetHomeProject() (models.Project, error) {
	if _, ok := ms.project[uuid.Max()]; !ok {
		return models.Project{}, store.ErrProjectNotFound
	}
	return ms.project[uuid.Max()].Project, nil
}

// CreateProject implements [models.Store].
func (ms *Store) CreateProject(projectName string) (models.Project, error) {
	project := models.Project{ProjecId: uuid.New(), ProjectName: projectName}
	ms.project[project.ProjecId] = home{Project: project, Table: models.Table{Rows: make([]models.TableRow, 0)}}
	return project, nil
}

// GetTable implements [models.Store].
func (ms *Store) GetTable(projectId uuid.UUID) (models.Table, error) {
	if _, ok := ms.project[projectId]; !ok {
		return models.Table{}, store.ErrProjectNotFound
	}
	return ms.project[projectId].Table, nil
}
