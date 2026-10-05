package routemodels

type URLQueryKey string

const (
	UrlQueryKeyOption URLQueryKey = "option" // ?option=
)

type DependencySuggestionType string

const (
	AddActivity DependencySuggestionType = "add_activity" // start to start relationship
)
