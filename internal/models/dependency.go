package models

type DependencyRelationship string

const (
	FS DependencyRelationship = "FS" // finish to start relationship
	SS DependencyRelationship = "SS" // start to start relationship
	SF DependencyRelationship = "SF" // start to finish relationship
	FF DependencyRelationship = "FF" // finish to finish relationship
)
