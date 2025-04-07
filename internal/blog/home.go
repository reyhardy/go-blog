package blog

import (
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

func Home(view string) gomponents.Node {
	return layout.Layout(
		header(view),
		body(view),
	)
}

func header(view string) gomponents.Node {
	return html.Header(
		// components.NavbarComponent(NavbarButtonId, element.ButtonElement("", "button", "Add Post", "")),
		gomponents.If(view == "post", NavbarAddPostButton()),
		gomponents.If(view == "form", NavbarBackButton()),
		html.HGroup(
			html.H1(gomponents.Text("this is go-blog")),
			html.P(gomponents.Text("we use picocss for styling")),
		),
	)
}

func body(view string) gomponents.Node {
	return html.Main(
		html.ID("main"),
		// html.Code(
		// 	html.Pre(
		// 		html.Data("text", "ctx.signals.JSON()"),
		// 	),
		// ),
		gomponents.If(view == "form", AddForm()),
		gomponents.If(view == "post", html.Div(
			html.ID(PostsId),
			html.Data("on-load", datastar.GetSSE("/posts")),
		)),
	)
}
