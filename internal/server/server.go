package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gtantech/p2/internal/models"
)

type server struct {
	store models.Store
}

func NewServer(config models.ServerConfig, store models.Store) *http.Server {
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
