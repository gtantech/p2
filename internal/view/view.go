package view

import "github.com/a-h/templ"

type View struct {
}

func (v *View) Home() templ.Component {
	return home()
}
