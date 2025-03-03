package blog

import (
	"github.com/reyhardy/go-blog/template/components"
	"github.com/reyhardy/go-blog/template/element"
	"github.com/reyhardy/go-blog/template/layout"
	datastar "github.com/starfederation/datastar/sdk/go"
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

const (
	NavbarButtonId string = "add-back-button"
	PostsId        string = "posts"
	FormId         string = "form"
)

func Home() gomponents.Node {
	return layout.Layout(
		header(),
		main(),
	)
}

func header() gomponents.Node {
	return html.Header(
		components.NavbarComponent(NavbarButtonId, element.ButtonElement("", "button", "No Name")),
		html.HGroup(
			html.H1(gomponents.Text("this is go-blog")),
			html.P(gomponents.Text("we use picocss for styling")),
		),
	)
}

func main() gomponents.Node {
	return html.Main(
		html.Code(
			html.Pre(
				html.Data("text", "ctx.signals.JSON()"),
			),
		),
		html.Data("signals", "{'view': 'posts'}"),
		html.Div(
			html.ID(PostsId),
			html.Data("on-load", datastar.GetSSE("/posts")),
			html.Data("show", "$view ==='posts'"),
		),
		html.Div(
			html.ID(FormId),
			html.Data("show", "$view ==='form'"),
		),
	)
}
