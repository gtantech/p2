package viewmodels

import (
	"github.com/a-h/templ"
	"github.com/gtantech/p2/internal/models/htmxmodels"
)

type HtmxInput struct {
	HtmlInput
	Hx
}

type HtmxButton struct {
	HtmlButton
	Hx
}

type Hx struct {
	htmxmodels.HxMethod
	htmxmodels.HxRequest
	htmxmodels.HxResponse
	htmxmodels.HxOptions
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
	add("hx-swap-oob", h.HxSwapOob)

	// Options
	add("hx-confirm", h.HxConfirm)
	add("hx-indicator", h.HxIndicator)

	return items
}
