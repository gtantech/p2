package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/store"
)

type Server struct {
	*http.Server
	port int
}

func NewServer(config models.ServerConfig) *Server {
	newServer := &Server{
		port: config.Port,
	}
	store := &store.MockStore{}
	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(store),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	newServer.Server = server

	return newServer
}
