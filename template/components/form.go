package components

import (
	"github.com/reyhardy/go-blog/template/element"
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func Form() gomponents.Node {
	return html.Form(
		gomponents.Attr("data-signals", "{input: '', title: 'Title', author: 'Author', content: 'Content'}"),
		gomponents.Attr("data-on-submit", "@post('/post', {contentType: 'form'})"),
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
	)
}
