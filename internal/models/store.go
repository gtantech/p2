package models

import "uuid"

type Store interface {
	CreateProject(projectName string) Project
	GetHomeProject() Project
	GetTable(projectId uuid.UUID) Table
}
