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
	port  int
	store *store.Store
}

func NewServer(config ServerConfig) *http.Server {
	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("failed to open database with error: %v", err)
	}
	NewServer := &Server{
		port:  config.Port,
		store: store.NewStoreFromDb(db.New(database)),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
