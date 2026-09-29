package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/gtantech/p2/internal/view"
	"github.com/gtantech/p2/static"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", templ.Handler(view.Home()).ServeHTTP)
	r.Get("/static/home_style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(static.StaticHomeCss)
	})
	http.ListenAndServe(":8080", r)
}
