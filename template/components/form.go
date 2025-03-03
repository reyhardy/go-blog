package components

import (
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func FormComponent(id, legend string, onSubmit, button gomponents.Node, input ...gomponents.Node) gomponents.Node {
	return html.Div(
		html.Form(
			html.ID(id),
			onSubmit,
			html.FieldSet(
				html.Legend(gomponents.Text(legend)),
				gomponents.Group(input),
				button,
			),
		),
	)
}
