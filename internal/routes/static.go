package routes

import (
	"net/http"

	"github.com/gtantech/p2/static"
)

type StaticRoutes struct {
}

func NewStaticRoutes() *StaticRoutes {
	return &StaticRoutes{}
}

func (rt *StaticRoutes) GetHomeStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(static.StaticHomeCss)
}
