package blog

import (
	"github.com/reyhardy/go-blog/template/components"
	"github.com/reyhardy/go-blog/template/layout"
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func Home() gomponents.Node {
	return layout.Layout(
		components.Navbar("Home", "About", "Contact"),
		html.HGroup(
			html.H1(gomponents.Text("this is go-blog")),
			html.P(gomponents.Text("we use picocss for styling")),
		),
		html.Div(
			html.Class("grid"),
			components.ModalForm("Add Post", "Add Post"),
			html.Div(
				gomponents.Attr("data-on-load", "@get('/post')"),
				html.ID("postlist"),
			),
		),
	)
}
