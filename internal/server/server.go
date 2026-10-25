package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gtantech/p2/internal/models"
)

type Server struct {
	*http.Server
	port  int
	store models.Store
}

func NewServer(config models.ServerConfig, store models.Store) *Server {
	newServer := &Server{
		port:  config.Port,
		store: store,
	}
	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	newServer.Server = server

	return newServer
}
