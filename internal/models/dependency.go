package models

import "uuid"

type Dependency struct {
	DependencyId          uuid.UUID
	SuccessorActivityId   uuid.UUID
	PredecessorActivityId uuid.UUID
}
