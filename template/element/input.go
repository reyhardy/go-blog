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
			html.Data(fmt.Sprintf("bind-input.%s", name), ""),
			html.ID(name),
			html.Name(name),
			html.Type(inputType),
			html.Value(value),
			// html.Placeholder(value),
			gomponents.Group(attr),
		),
	)
}

func Textarea(label, name, value string, rows int, attr ...gomponents.Node) gomponents.Node {
	return html.Label(
		gomponents.Text(label),
		html.For(name),
		html.Textarea(
			html.Data(fmt.Sprintf("bind-input.%s", name), ""),
			html.ID(name),
			html.Name(name),
			html.Value(value),
			// html.Placeholder(value),
			html.Rows(fmt.Sprintf("%d", rows)),
			gomponents.Group(attr),
		),
	)
}
