package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gtantech/p2/internal/routes"
)

func (s *Server) RegisterRoutes(routes *routes.Routes) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", routes.HelloWorldHandler)

	r.Route("/web", func(r chi.Router) {
		r.Get("/dependency/add", routes.DisplayDependenciesToAdd)
	})

	r.Route("/component", func(r chi.Router) {
		r.Post("/table/row/empty", routes.PostEmptyTableRow)
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/table", func(r chi.Router) {
			r.Put("/row/activity/{id}/name", routes.PutActivityNameUpdateFromTableHandler)
			r.Put("/row/activity/{id}/dependency", routes.PutActivityDependencyUpdateFromTableHandler)
			r.Put("/row/activity/{id}/duration", routes.PutActivityDurationUpdateFromTableHandler)
		})
	})

	r.Get("/static/home_style.css", routes.GetHomeStyle)
	return r
}
