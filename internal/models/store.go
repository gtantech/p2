package models

import "uuid"

type Store interface {
	CreateProject(projectName string) Project
	GetTable(projectId uuid.UUID) Table
}
