package viewmodels

import (
	"github.com/a-h/templ"
)

type HxMethod struct {
	HxGet    string
	HxPost   string
	HxPut    string
	HxPatch  string
	HxDelete string
}

type HxRequest struct {
	HxTrigger string
	HxVals    string
	HxInclude string
	HxHeaders string
}

type HxResponse struct {
	HxTarget string
	HxSwap   string
}

type HxOptions struct {
	HxConfirm   string
	HxIndicator string
}

type HtmxInput struct {
	HtmlInput
	Hx
}

type Hx struct {
	HxMethod
	HxRequest
	HxResponse
	HxOptions
}

func (h Hx) Items() []templ.KeyValue[string, any] {
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

	// Options
	add("hx-confirm", h.HxConfirm)
	add("hx-indicator", h.HxIndicator)

	return items
}
