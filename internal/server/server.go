package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/database/sqlitedb"
	"github.com/gtantech/p2/internal/models/servermodels"
	"github.com/gtantech/p2/internal/routes"
	"github.com/gtantech/p2/internal/view"
)

type Server struct {
	*http.Server
	port int
}

func NewServer(config servermodels.ServerConfig) *Server {
	NewServer := &Server{
		port: config.Port,
	}

	sqliteQueries := sqlitedb.New(sqlitedb.NewSQLiteDb(":memory:"))
	if os.Getenv("env") == "test" {
		insertMockData(sqliteQueries)
	}
	// Declare Server config
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(
			routes.NewViewRoutes(&view.HttpViewTemplAdapter{}, sqlitedb.NewStoreViewSqliteAdapter(sqliteQueries), sqlitedb.NewStoreFormSqliteAdapter(sqliteQueries)),
			&routes.StaticRoutes{},
			routes.NewFormRoutes(sqlitedb.NewStoreFormSqliteAdapter(sqliteQueries))),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	NewServer.Server = server

	return NewServer
}

func insertMockData(sqliteQueries *sqlitedb.Queries) {
	ctx := context.Background()
	projectId := uuid.Max()

	sqliteQueries.InsertProject(ctx, sqlitedb.InsertProjectParams{ID: projectId.String(), DispName: "mock project"})

	defaultSortRankStep := int64(1000)
	rows := []sqlitedb.Activity{
		{ID: uuid.New().String(), ProjectID: projectId.String(), DispName: "A", Duration: int64(1 * time.Minute)},
		{ID: uuid.New().String(), ProjectID: projectId.String(), DispName: "B", Duration: int64(2 * time.Minute)},
		{ID: uuid.New().String(), ProjectID: projectId.String(), DispName: "C", Duration: int64(3 * time.Minute)},
	}

	for i, row := range rows {
		sqliteQueries.InsertActivity(ctx, sqlitedb.InsertActivityParams{ID: row.ID, ProjectID: row.ProjectID, DispName: row.DispName, Duration: row.Duration})
		sqliteQueries.InsertActivityOrdering(ctx, sqlitedb.InsertActivityOrderingParams{
			ID:                  uuid.New().String(),
			ProjectID:           row.ProjectID,
			SuccessorActivityID: row.ID,
			SortRank:            int64(i) * defaultSortRankStep,
		})
	}

	sqliteQueries.InsertDependency(ctx, sqlitedb.InsertDependencyParams{ID: uuid.New().String(), ProjectID: projectId.String(), Relationship: "FS", PredecessorActivityID: rows[0].ID, SuccessorActivityID: rows[1].ID})
	sqliteQueries.InsertDependency(ctx, sqlitedb.InsertDependencyParams{ID: uuid.New().String(), ProjectID: projectId.String(), Relationship: "FS", PredecessorActivityID: rows[1].ID, SuccessorActivityID: rows[2].ID})
	sqliteQueries.InsertDependency(ctx, sqlitedb.InsertDependencyParams{ID: uuid.New().String(), ProjectID: projectId.String(), Relationship: "FS", PredecessorActivityID: rows[0].ID, SuccessorActivityID: rows[2].ID})
}
