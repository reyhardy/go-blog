package element

import (
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func InputElement(label, name, inputType, dataBind string, attr ...gomponents.Node) gomponents.Node {
	return html.Label(
		gomponents.Text(label),
		html.For(name),
		html.Input(
			html.ID(name),
			html.Name(name),
			html.Type(inputType),
			html.Data("bind", dataBind),
			gomponents.Group(attr),
		),
	)
}

func Textarea(label, name, dataBind string, attr ...gomponents.Node) gomponents.Node {
	return html.Label(
		gomponents.Text(label),
		html.For(name),
		html.Textarea(
			html.ID(name),
			html.Name(name),
			html.Data("bind", dataBind),
			gomponents.Group(attr),
		),
	)
}
