package element

import (
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func ButtonElement(id, btnType, btnText string, attr ...gomponents.Node) gomponents.Group {
	return gomponents.Group{
		html.Button(
			html.ID(id),
			html.Type(btnType),
			gomponents.Text(btnText),
			gomponents.Group(attr),
		),
	}
}
