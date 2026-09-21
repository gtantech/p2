package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gtantech/p2/internal/routes"
	"github.com/gtantech/p2/static"
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

	r.Handle("/static/*", http.StripPrefix(
		"/static/",
		http.FileServer(http.FS(static.StaticHomeCss)),
	))
	return r
}
