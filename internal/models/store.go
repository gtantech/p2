package models

import "uuid"

type Store interface {
	// Creates a new project
	CreateProject(projectName string) (Project, error)

	// Gets the home project
	// GetHomeProject returns ErrProjectNotFound if no project exists.
	GetHomeProject() (Project, error)

	// Gets the table with rows associated with the projectId sorted by sort rank.
	// GetTable returns ErrProjectNotFound if projectId does not exist.
	GetTable(projectId uuid.UUID) (Table, error)
}
