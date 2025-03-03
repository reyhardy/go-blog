package layout

import (
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/components"
	"maragu.dev/gomponents/html"
)

func Layout(children ...gomponents.Node) gomponents.Node {
	return components.HTML5(components.HTML5Props{
		Title: "go-blog",
		Head: []gomponents.Node{
			// html.Script(html.Type("module"), html.Src("/static/js/datastar-1-0-0-beta-1-709a1d5426cfe2c0.js")),
			html.Script(html.Type("module"), html.Src("/static/js/datastar-1-0-0-beta-8-6ebd8eefd5e077e9.js")),
			html.Link(html.Rel("stylesheet"), html.Href("/static/css/pico.min.css")),
		},
		Body: []gomponents.Node{
			html.Body(
				html.Class("container"),
				gomponents.Group(children),
			),
		},
	})
}
