package components

import (
	"fmt"

	"github.com/reyhardy/go-blog/template/element"
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func ModalForm(id, btnName, modalHeader string, component ...gomponents.Node) gomponents.Node {
	return html.Div(
		html.Dialog(
			html.ID(fmt.Sprintf("dialog-%s", id)),
			html.Data("ref", "dialog"),
			html.Article(
				// html.Data("on-click", "!$dialog.open ? null : $dialog.close()"),
				html.Header(
					element.ButtonElement("", "", "",
						html.Aria("label", "Close"),
						html.Rel("prev"),
						html.Data("on-click", "$dialog.close()"),
					),
					html.H1(html.Strong(gomponents.Text(modalHeader))),
				),
				gomponents.Group(component),
			),
		),
		element.ButtonElement("", "button", btnName, html.Data("on-click", "$dialog.showModal(); console.log($dialog)")),
	)
}
