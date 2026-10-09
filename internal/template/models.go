package template

import (
	"uuid"

	"github.com/a-h/templ"
)

type InputAutocomplete string

const (
	On  InputAutocomplete = "on"
	Off InputAutocomplete = "off"
)

type htmlInput struct {
	Id           string
	Name         string
	Class        string
	Placeholder  string
	Autocomplete InputAutocomplete
	Value        string
}

type htmlDiv struct {
	Id    string
	Class string
}

type htmlButton struct {
	Class string
	Title string
	Text  string
}

type htmxButton struct {
	htmlButton
	hx
}

type dependencyTableRow struct {
	Id                         string
	Class                      string
	ActivityTextInputParams    *htmlInput
	DependencyWrapperDivParams *htmlDiv
	DependencyTextInputParams  *htmlInput
	DurationTextInputParams    *htmlInput
	TableRowAddBtnParams       *htmxButton
}

type dependencyTable struct {
	ProjectId uuid.UUID
	Rows      []dependencyTableRow
}

type homeParams struct {
	Table                 dependencyTable
	TableWrapperDivParams *htmlDiv
	ContainerDivParams    *htmlDiv
	PageTitle             string
}

type hxMethod struct {
	HxGet    string
	HxPost   string
	HxPut    string
	HxPatch  string
	HxDelete string
}

type hxRequest struct {
	HxTrigger string
	HxVals    string
	HxInclude string
	HxHeaders string
}

type hxResponse struct {
	HxTarget  string
	HxSwap    string
	HxSwapOob string
}

type hxOptions struct {
	HxConfirm   string
	HxIndicator string
}

type hx struct {
	hxMethod
	hxRequest
	hxResponse
	hxOptions
}

func (h hx) Items() []templ.KeyValue[string, any] {
	items := make([]templ.KeyValue[string, any], 0, 13)

	add := func(key, value string) {
		if value != "" {
			items = append(items, templ.KV[string, any](key, value))
		}
	}

	// Methods
	add("hx-get", h.HxGet)
	add("hx-post", h.HxPost)
	add("hx-put", h.HxPut)
	add("hx-patch", h.HxPatch)
	add("hx-delete", h.HxDelete)

	// Request
	add("hx-trigger", h.HxTrigger)
	add("hx-vals", h.HxVals)
	add("hx-include", h.HxInclude)
	add("hx-headers", h.HxHeaders)

	// Response
	add("hx-target", h.HxTarget)
	add("hx-swap", h.HxSwap)
	add("hx-swap-oob", h.HxSwapOob)

	// Options
	add("hx-confirm", h.HxConfirm)
	add("hx-indicator", h.HxIndicator)

	return items
}
