package routes

import (
	"context"
	"uuid"

	"github.com/gtantech/p2/internal/models"
)

type Store interface {
	// Creates a new project
	CreateProject(ctx context.Context, projectName string) (models.Project, error)

	// Gets the home project
	// GetHomeProject returns ErrProjectNotFound if no project exists.
	GetHomeProject(ctx context.Context) (models.Project, error)

	// Gets the table with rows associated with the projectId sorted by sort rank.
	// GetTable returns ErrProjectNotFound if projectId does not exist.
	GetTable(ctx context.Context, projectId uuid.UUID) (models.Table, error)

	// CreateRowBelow adds a new empty row below the specified rowId.
	// CreateRowBelow returns ErrRowNotFound if no row with the rowId exists.
	CreateRowBelow(ctx context.Context, rowId uuid.UUID) (models.TableRow, error)
}
