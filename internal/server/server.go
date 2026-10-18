package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gtantech/p2/internal/database/sqlitedb"
	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/store"
	"github.com/gtantech/p2/internal/template"
)

type Server struct {
	*http.Server
	port int
}

func NewServer(config models.ServerConfig) *Server {
	newServer := &Server{
		port: config.Port,
	}
	t := template.NewTemplPresenter()
	var storeModel models.Model
	if os.Getenv("ENV") == "mock" {
		storeModel = store.NewSqliteStore(sqlitedb.New(sqlitedb.NewSQLiteDbWithMockValues(":memory:")))
	} else {
		storeModel = store.NewSqliteStore(sqlitedb.New(sqlitedb.NewSQLiteDb("./data.db")))
	}
	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(t, storeModel),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	newServer.Server = server

	return newServer
}
