package htmxmodels

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
