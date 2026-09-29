package server

import (
	"fmt"
	"net/http"
	"time"

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

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(routes.NewViewRoutes(&view.HttpViewTemplAdapter{}), &routes.StaticRoutes{}),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	NewServer.Server = server

	return NewServer
}
