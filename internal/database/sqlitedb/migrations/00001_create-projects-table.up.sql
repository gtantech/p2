-- projects definition

CREATE TABLE IF NOT EXISTS projects (
	id TEXT NOT NULL, disp_name TEXT NOT NULL,
	CONSTRAINT projects_pk PRIMARY KEY (id)
);