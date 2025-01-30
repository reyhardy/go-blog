package element

import (
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func ButtonElement(btnType, btnText string, attr ...gomponents.Node) gomponents.Node {
	return html.Button(
		html.Type(btnType),
		gomponents.Text(btnText),
		gomponents.Group(attr),
	)
}
