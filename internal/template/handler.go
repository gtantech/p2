package template

import "net/http"

func NewTemplateErrHandler(err string, httpStatus int) *templateErrHttpHandler {
	return &templateErrHttpHandler{error: err, status: httpStatus}
}

type templateErrHttpHandler struct {
	status int
	error  string
}

func (t *templateErrHttpHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	http.Error(w, t.error, t.status)
}
