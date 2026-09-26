package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gtantech/p2/internal/db"
	"github.com/gtantech/p2/internal/routes"
	"github.com/gtantech/p2/internal/store"
	_ "github.com/joho/godotenv/autoload"
)

type ServerConfig struct {
	Port int
}

type Server struct {
	*http.Server
	database *sql.DB
	port     int
}

func (s *Server) CloseDb() error {
	return s.database.Close()
}

func NewServer(config ServerConfig) *Server {
	database, err := db.NewSQLiteStorage(":memory:")
	if err != nil {
		log.Fatalf("failed to open database with error: %v", err)
	}

	store := store.NewStoreFromDb(db.New(database))

	NewServer := &Server{
		database: database,
		port:     config.Port,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(routes.NewRoutes(routes.NewStoreAdapter(store))),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	NewServer.Server = server

	return NewServer
}
