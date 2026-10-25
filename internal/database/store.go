package database

type DBQueries interface {
	//Close closes the unexported *sql.DB field
	Close() error
}
