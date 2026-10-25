package mock

import (
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models"
)

var _ models.Store = (*Store)(nil) //ensures MockStore implements models.Store at compile time

type Store struct {
}

// GetTable implements [models.Store].
func (ms *Store) GetTable(projectId uuid.UUID) models.Table {
	return models.Table{
		Rows: []models.TableRow{
			{RowId: uuid.MustParse("dcce6c98-b31b-4ed6-a2ed-15d24ae96b41"), ActivityName: "A", DependencyNames: []string{}, Duration: 5 * time.Minute},
			{RowId: uuid.MustParse("a4518b28-e597-4295-a532-3a501a75ab2b"), ActivityName: "B", DependencyNames: []string{}, Duration: 2 * time.Minute},
			{RowId: uuid.MustParse("aa5a833c-2407-464c-93c4-fa9edba03b52"), ActivityName: "C", DependencyNames: []string{"A", "B"}, Duration: 3 * time.Minute},
		},
	}
}

// GetHome implements [models.Store].
func (ms *Store) GetHome() models.Home {
	return models.Home{Name: "World"}
}
