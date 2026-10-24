package models

import (
	"net/http"
)

type Presenter interface {
	RegisterModel(model Model)
	DisplayHomeHandler() http.Handler
}

type Model interface {
	GetHome() Home
}
