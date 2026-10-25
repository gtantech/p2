package models

import "uuid"

type Store interface {
	GetTable(projectId uuid.UUID) Table
}
