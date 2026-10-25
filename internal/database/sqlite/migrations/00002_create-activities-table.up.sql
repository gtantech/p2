-- activities definition

CREATE TABLE IF NOT EXISTS activities (
	id TEXT NOT NULL,
	project_id TEXT NOT NULL,
	disp_name TEXT NOT NULL,
	duration INTEGER NOT NULL,
	CONSTRAINT activities_pk PRIMARY KEY (id),
	CONSTRAINT activities_projects_FK FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE ON UPDATE CASCADE
);