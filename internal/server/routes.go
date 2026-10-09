package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/routes"
)

func (s *Server) RegisterRoutes(presenter models.Presenter, store models.Model) http.Handler {
	rts := routes.NewRoutes(presenter, store)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", rts.HomeHandler)
	r.Get("/static/home_style.css", rts.GetHomeStyleHandler)
	r.Post("/form/json/table/row/empty/component", rts.PostFromRowPlusBtnReturnsEmptyTableRowHandler)
	r.Put("/form/json/table/row/activity/name", rts.PutActivityNameHandler)
	r.Put("/form/json/table/row/activity/duration", rts.PutActivityDurationHandler)
	return r
}
