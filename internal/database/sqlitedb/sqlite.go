package sqlitedb

import (
	"database/sql"
	"embed"
	"errors"
	"log"
	"time"
	"uuid"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/gtantech/p2/internal/models"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func NewSQLiteDb(dataSourceName string) *sql.DB {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Verify the connection is working
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatalf("Failed to create driver for sqlite database: %v", err)
	}
	if err := migrateDB(driver); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	return db
}

func NewSQLiteDbWithMockValues(dataSourceName string) *sql.DB {
	return InjectMockValues(NewSQLiteDb(dataSourceName))
}

func InjectMockValues(db *sql.DB) *sql.DB {
	projectId := uuid.Max()
	db.Exec("INSERT INTO projects (id, disp_name) VALUES (?, ?)", projectId, "mock project")

	activities := []models.Activity{
		{ActivityId: uuid.MustParse("dcce6c98-b31b-4ed6-a2ed-15d24ae96b41"), DisplayName: "A", Duration: 5 * time.Minute},
		{ActivityId: uuid.MustParse("a4518b28-e597-4295-a532-3a501a75ab2b"), DisplayName: "B", Duration: 4 * time.Minute},
		{ActivityId: uuid.MustParse("aa5a833c-2407-464c-93c4-fa9edba03b52"), DisplayName: "C", Duration: 3 * time.Minute},
	}
	for i, a := range activities {
		db.Exec("INSERT INTO activities (id, project_id, disp_name, duration) VALUES (?, ?, ?, ?)",
			a.ActivityId, projectId, a.DisplayName, a.Duration)
		db.Exec("INSERT INTO activities_ordering (id, project_id, successor_activity_id, sort_rank) VALUES (?, ?, ?, ?)",
			uuid.NewV7(), projectId, a.ActivityId, int64(i)*models.TableRowSortRankStep)
	}

	db.Exec("INSERT INTO dependencies (id, project_id, relationship, predecessor_activity_id, successor_activity_id) VALUES (?, ?, ?, ?, ?)",
		uuid.NewV7(), projectId, models.FS, activities[0].ActivityId, activities[1].ActivityId)
	db.Exec("INSERT INTO dependencies (id, project_id, relationship, predecessor_activity_id, successor_activity_id) VALUES (?, ?, ?, ?, ?)",
		uuid.NewV7(), projectId, models.FS, activities[0].ActivityId, activities[2].ActivityId)
	db.Exec("INSERT INTO dependencies (id, project_id, relationship, predecessor_activity_id, successor_activity_id) VALUES (?, ?, ?, ?, ?)",
		uuid.NewV7(), projectId, models.FS, activities[1].ActivityId, activities[2].ActivityId)
	return db
}

func migrateDB(driver database.Driver) error {
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		source,
		"main",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
