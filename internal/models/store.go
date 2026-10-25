package models

import (
	"context"
	"uuid"
)

type Store interface {
	// Creates a new project
	CreateProject(ctx context.Context, projectName string) (Project, error)

	// Gets the home project
	// GetHomeProject returns ErrProjectNotFound if no project exists.
	GetHomeProject(ctx context.Context) (Project, error)

	// Gets the table with rows associated with the projectId sorted by sort rank.
	// GetTable returns ErrProjectNotFound if projectId does not exist.
	GetTable(ctx context.Context, projectId uuid.UUID) (Table, error)
}
