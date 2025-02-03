package components

import (
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func Navbar(navbar ...string) gomponents.Node {
	// navbar := []string{"Home", "About", "Contact"}
	return html.Nav(
		html.Ul(
			html.Li(html.Strong(gomponents.Text("Go-Blog"))),
		),
		html.Ul(
			gomponents.Map(navbar, func(nb string) gomponents.Node {
				return html.Li(html.A(html.Href("#"), gomponents.Text(nb)))
			}),
		),
	)
}
