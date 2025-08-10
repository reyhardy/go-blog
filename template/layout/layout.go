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
			html.Script(html.Type("module"), html.Src("/static/js/datastar.js")),
			html.Link(html.Rel("stylesheet"), html.Href("/static/css/pico.min.css")),
		},
		Body: []gomponents.Node{
			gomponents.Group(children),
		},
		HTMLAttrs: []gomponents.Node{},
	})
}
