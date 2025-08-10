package blog

import (
	"github.com/reyhardy/go-blog/template/layout"
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func Home() gomponents.Node {
	return layout.Layout(
		header(),
		main(),
		footer(),
	)
}

func header() gomponents.Node {
	return html.Header(
		html.Class("container-fluid"),
		html.Nav(
			html.Ul(
				html.Li(
					html.HGroup(
						html.H2(
							gomponents.Text("Go Blog"),
						),
						html.P(
							gomponents.Text("A simple blog application built with Go."),
						),
					),
				),
			),
			html.Ul(
				html.Li(
					html.ID("navbar-button"),
					ModalAddPost(),
				),
			),
		),
	)
}

func main() gomponents.Node {
	return html.Main(
		html.Class("container-fluid"),
		html.Div(
			html.ID("posts"),
			html.Data("on-load", `@get('/api/posts')`),
		),
	)
}

func footer() gomponents.Node {
	return html.Footer(
		html.Class("container-fluid"),
		html.P(
			gomponents.Text("I'm a Footer"),
		),
	)
}
