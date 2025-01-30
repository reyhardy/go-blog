package element

import (
	"fmt"

	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func InputElement(label, name, value, inputType string, attr ...gomponents.Node) gomponents.Node {
	return html.Label(
		gomponents.Text(label),
		html.For(name),
		html.Input(
			gomponents.Attr(fmt.Sprintf("data-bind-%s", name)),
			html.ID(name),
			html.Name(name),
			html.Type(inputType),
			html.Value(value),
			gomponents.Group(attr),
		),
	)
}

func Textarea(label, name, value string, attr ...gomponents.Node) gomponents.Node {
	return html.Label(
		gomponents.Text(label),
		html.For(name),
		html.Textarea(
			gomponents.Attr(fmt.Sprintf("data-bind-%s", name)),
			html.ID(name),
			html.Name(name),
			html.Value(value),
			gomponents.Group(attr),
		),
	)
}
