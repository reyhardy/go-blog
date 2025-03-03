package components

import (
	"fmt"

	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func CardComponent(id, title string, footer gomponents.Group, children ...gomponents.Node) gomponents.Node {
	return html.Div(
		html.ID(fmt.Sprintf("post-%s", id)),
		html.Article(
			html.Header(html.H1(gomponents.Text(title))),
			gomponents.Group(children),
			html.Footer(
				footer,
			),
		),
	)

}
