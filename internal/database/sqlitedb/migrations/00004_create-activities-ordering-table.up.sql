-- activities_ordering definition

CREATE TABLE IF NOT EXISTS activities_ordering (
	id TEXT NOT NULL,
	project_id TEXT NOT NULL,
	successor_activity_id TEXT NOT NULL,
	sort_rank INTEGER NOT NULL,
	CONSTRAINT activities_ordering_pk PRIMARY KEY (id),
	CONSTRAINT activities_ordering_successor_activity_UK UNIQUE (successor_activity_id),
	CONSTRAINT activities_ordering_activities_successor_FK FOREIGN KEY (successor_activity_id) REFERENCES activities(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT activities_ordering_projects_FK FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE ON UPDATE CASCADE
);