package components

import (
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func NavbarComponent(id string, children ...gomponents.Node) gomponents.Node {
	return html.Nav(
		html.ID(id),
		html.Ul(
			html.Li(html.Strong(gomponents.Text("Go-Blog"))),
		),
		html.Ul(
			html.Li(
				gomponents.Group(children),
			),
		),
	)
}
