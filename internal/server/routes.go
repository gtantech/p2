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
		r.Get("/table/dependency/row/empty", routes.GetEmptyTableRow)
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/table", func(r chi.Router) {
			r.Post("/row/activity/{id}/name", routes.PostActivityNameUpdateFromTableHandler)
			r.Post("/row/activity/{id}/dependency", routes.PostActivityDependencyUpdateFromTableHandler)
			r.Post("/row/activity/{id}/duration", routes.PostActivityDurationUpdateFromTableHandler)
		})
	})

	r.Get("/static/home_style.css", routes.GetHomeStyle)
	return r
}
