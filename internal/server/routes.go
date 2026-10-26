package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gtantech/p2/internal/routes"
)

func (s *server) RegisterRoutes() http.Handler {
	rts := routes.New(s.store)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", rts.GetHome)
	r.Get("/static/home_style.css", rts.GetHomeStyle)

	r.Route("/form", func(r chi.Router) {
		r.Route("/json", func(r chi.Router) {
			r.Route("/table", func(r chi.Router) {
				r.Route("/row", func(r chi.Router) {
					r.Post("/empty/below", rts.PostFromRowPlusBtnReturnsEmptyTableRow)
				})
			})
		})
	})

	return r
}
