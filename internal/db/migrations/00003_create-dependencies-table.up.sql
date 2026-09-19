-- dependencies definition

CREATE TABLE IF NOT EXISTS dependencies (
	id TEXT NOT NULL,
	project_id TEXT NOT NULL,
	relationship TEXT NOT NULL,
	predecessor_activity_id TEXT NOT NULL,
	successor_activity_id TEXT NOT NULL,
	CONSTRAINT dependencies_pk PRIMARY KEY (id),
	CONSTRAINT dependencies_activities_predecessor_FK FOREIGN KEY (predecessor_activity_id) REFERENCES activities(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT dependencies_activities_successor_FK FOREIGN KEY (successor_activity_id) REFERENCES activities(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT dependencies_projects_FK FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE ON UPDATE CASCADE
);