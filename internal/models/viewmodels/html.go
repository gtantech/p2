package viewmodels

type HtmlInput struct {
	Id           string
	Name         string
	Class        string
	Placeholder  string
	Autocomplete InputAutocomplete
	Value        string
}

type HtmlButton struct {
	Class string
	Title string
	Text  string
}

type HtmlDependencyTableRowButton struct {
	HtmlButton
	HxPost string
	HxVals string
}
