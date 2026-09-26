package models

import (
	"uuid"
)

type StoreProject struct {
	ID          uuid.UUID
	DisplayName string
}

type StoreCreateProjectParams struct {
	DisplayName string
}

type StoreUpdateProjectParams struct {
	Id          uuid.UUID
	DisplayName string
}

type Dependency struct {
	ID                    uuid.UUID
	ProjectID             uuid.UUID
	Relationship          RelationshipType
	PredecessorActivityID uuid.UUID
	SuccessorActivityID   uuid.UUID
}

type GetPredecessorNamesBySuccessorResult struct {
	DependencyID            uuid.UUID
	Relationship            RelationshipType
	PredecessorActivityID   uuid.UUID
	PredecessorActivityName string
}

type CreateDepdencencyParams struct {
	ProjectID             uuid.UUID
	Relationship          RelationshipType
	PredecessorActivityID uuid.UUID
	SuccessorActivityID   uuid.UUID
}

type UpdateDepdencencyParams struct {
	ID                    uuid.UUID
	Relationship          RelationshipType
	PredecessorActivityID uuid.UUID
	SuccessorActivityID   uuid.UUID
}
