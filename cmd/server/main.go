package main

import (
	"fmt"
	"net/http"

	"github.com/gtantech/p2/internal/routes"
	"github.com/gtantech/p2/internal/store"
	"github.com/gtantech/p2/internal/template"
)

func main() {
	t := template.NewTemplPresenter()
	store := &store.MockStore{}
	r := routes.NewRoutes(t, store)
	http.HandleFunc("/", r.HomeHandler)

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
