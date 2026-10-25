package models

import "uuid"

type Store interface {
	GetHome() Home
	GetTable(projectId uuid.UUID) Table
}
