package components

import (
	"github.com/reyhardy/go-blog/template/element"
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func ModalForm(btnName, modalHeader string) gomponents.Node {
	return html.Div(
		gomponents.Attr("data-signals", "{open: false}"),
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
					gomponents.Attr("data-on-submit", "@post('/post', {contentType: 'form'}); $open = false; @setAll('input.', '')"),
					html.FieldSet(
						element.InputElement("Title", "title", "", "text", gomponents.Attr("data-bind", "input.title")),
						element.InputElement("Author", "author", "", "text", gomponents.Attr("data-bind", "input.author")),
						element.Textarea("Content", "content", "", html.Rows("10"), gomponents.Attr("data-bind", "input.content")),
						element.ButtonElement("submit", "Submit"),
					),
				),
			),
		),
	)
}
