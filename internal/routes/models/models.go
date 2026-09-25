package models

import "uuid"

type RelationshipType string

const (
	SS RelationshipType = "SS" // start to start relationship
	FS RelationshipType = "FS" // finish to start relationship
	SF RelationshipType = "SF" // start to finish relationship
	FF RelationshipType = "FF" // finish to finish relationship
)

type PostActivityDependencyUpdateFromTableParams struct {
	Relationship RelationshipType `json:"relationship"`
}

type PostEmptyTableRow struct {
	ProjectId uuid.UUID `json:"projectId"`
}
