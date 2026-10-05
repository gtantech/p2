package storemodels

type Activity struct {
	ID        string
	ProjectID string
	DispName  string
	Duration  int64
}

type ActivityNameWithId struct {
	ID       string
	DispName string
}
