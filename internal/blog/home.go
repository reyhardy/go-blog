package blog

import (
	"fmt"

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
		body(),
	)
}

func header(view string) gomponents.Node {
	return html.Header(
		gomponents.If(view == "posts", NavbarAddPostButton()),
		gomponents.If(view == "form", NavbarBackButton()),
		html.HGroup(
			html.H1(gomponents.Text("this is go-blog")),
			html.P(gomponents.Text("we use picocss for styling")),
		),
	)
}

func body() gomponents.Node {
	return html.Main(
		html.ID("main"),
		html.Code(
			html.Pre(
				html.Data("text", "ctx.signals.JSON()"),
			),
		),
		html.Data("signals", fmt.Sprintf("{'view': '%s'}", PostsId)),
		html.Div(
			html.ID(PostsId),
			html.Data("show", fmt.Sprintf("$view === '%s'", PostsId)),
			html.Data("on-load", datastar.GetSSE("/f/posts")),
		),
		html.Div(
			html.ID(FormId),
			html.Data("show", fmt.Sprintf("$view === '%s'", FormId)),
		),
	)
}
