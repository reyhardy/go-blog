package components

import (
	"github.com/reyhardy/go-blog/template/element"
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func ModalForm(btnName, modalHeader string) gomponents.Node {
	return html.Div(
		gomponents.Attr("data-signals", "{open: false, input: ''}"),
		element.ButtonElement("button", btnName, gomponents.Attr("data-on-click", "$open = true")),
		html.Dialog(
			html.ID("dialog"),
			gomponents.Attr("data-attr-open", "$open"),
			gomponents.Attr("data-on-keydown__window", "evt.key === 'Escape' ? $open = false : null"),
			html.Article(
				gomponents.Attr("data-on-click__outside__capture", "$open ? $open = false : null"),
				html.Header(
					element.ButtonElement("", "", html.Aria("label", "Close"), html.Rel("prev"), gomponents.Attr("data-on-click", "$open = !$open")),
					html.H1(html.Strong(gomponents.Text(modalHeader))),
				),
				html.Form(
					gomponents.Attr("data-on-submit", "@post('/addpost', {contentType: 'form'}); $open = false; $input = ''"),
					html.FieldSet(
						element.InputElement("Title", "title", "", "text"),
						element.InputElement("Author", "author", "", "text"),
						element.Textarea("Content", "content", "", html.Rows("10")),
						html.Div(
							html.Class("grid"),
							element.ButtonElement("submit", "Submit"),
							element.ButtonElement("reset", "Reset", gomponents.Attr("data-on-click", "$input=''")),
						),
					),
				),
			),
		),
	)
}
