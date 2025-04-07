package element

import (
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

const (
	BtnNeutral = "btn-neutral"
	BtnPrimary = "btn-primary"
	BtnError   = "btn-error"
	BtnSuccess = "btn-success"
	BtnWarning = "btn-warning"
)

func ButtonElement(id, btnType, btnText, btnClass string, attr ...gomponents.Node) gomponents.Node {
	switch btnClass {
	case BtnWarning:
		return gomponents.Group{
			html.Button(
				html.ID(id),
				html.Style("--pico-background-color: rgb(217, 53, 38); --pico-border-color: rgb(217, 53, 38);"),
				html.Type(btnType),
				gomponents.Text(btnText),
				gomponents.Group(attr),
			),
		}
	case BtnSuccess:
		return gomponents.Group{
			html.Button(
				html.ID(id),
				html.Style("--pico-background-color: rgb(255, 193, 7); --pico-border-color: rgb(255, 193, 7);"),
				html.Type(btnType),
				gomponents.Text(btnText),
				gomponents.Group(attr),
			),
		}
	default:
		return gomponents.Group{
			html.Button(
				html.ID(id),
				html.Type(btnType),
				gomponents.Text(btnText),
				gomponents.Group(attr),
			),
		}
	}
}
