package template

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/gtantech/p2/internal/models"
)

type TemplPresenter struct {
	model models.Model
}

func NewTemplPresenter() *TemplPresenter {
	return &TemplPresenter{}
}

// RegisterModel implements [models.Presenter].
func (t *TemplPresenter) RegisterModel(model models.Model) {
	t.model = model
}

// DisplayHomeHandler implements [models.Presenter].
func (t *TemplPresenter) DisplayHomeHandler() http.Handler {
	params := t.model.GetHome()
	component := Hello(params.Name)
	return templ.Handler(component)
}

var _ models.Presenter = (*TemplPresenter)(nil) //ensures TemplPresenter implements models.Presenter at compile time
