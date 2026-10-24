package models

type Home struct {
	Name string
}

type DisplayHomeParams struct {
	Home
	HttpParams
}
