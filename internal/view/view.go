package view

import "github.com/a-h/templ"

type View struct {
}

func NewView() *View {
	return &View{}
}

func (v *View) Home() templ.Component {
	return home()
}
