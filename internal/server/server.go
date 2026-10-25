package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/routes"
)

type server struct {
	store routes.Store
}

func NewServer(config models.ServerConfig, store routes.Store) *http.Server {
	newServer := &server{
		store: store,
	}
	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
