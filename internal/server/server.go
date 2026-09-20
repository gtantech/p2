package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gtantech/p2/internal/db"
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
	store    *store.Store
}

func (s *Server) CloseDb() error {
	return s.database.Close()
}

func NewServer(config ServerConfig) *Server {
	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("failed to open database with error: %v", err)
	}
	NewServer := &Server{
		database: database,
		port:     config.Port,
		store:    store.NewStoreFromDb(db.New(database)),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	NewServer.Server = server

	return NewServer
}
