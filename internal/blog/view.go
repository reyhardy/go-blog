package blog

import (
	"fmt"

	"github.com/reyhardy/go-blog/template/components"
	"github.com/reyhardy/go-blog/template/element"
	datastar "github.com/starfederation/datastar/sdk/go"
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func PostList(postList Posts) gomponents.Node {
	return html.Div(
		html.ID(PostsId),
		gomponents.Map(postList, func(post *Post) gomponents.Node {
			return PostCard(post)
		}),
	)
}

func PostCard(post *Post) gomponents.Node {
	return html.Div(
		html.ID(fmt.Sprintf("post-%s", post.ID)),
		html.Article(
			html.Header(html.H1(gomponents.Text(post.Title))),
			html.P(gomponents.Text(post.Content)),
			html.P(html.Cite(gomponents.Text(fmt.Sprintf("- %s", post.Author)))),
			// html.Footer(
			// 	html.P(gomponents.Text(fmt.Sprintf("Created at: %s", post.UpdatedAt.Local().Format("02 Jan 2006 @ 03:04 pm")))),
			// ),
			html.Div(
				html.Class("grid"),
				element.ButtonElement("delete-btn", "button", "Delete", element.BtnWarning, html.Data("on-click", datastar.DeleteSSE("/post/%s", post.ID))),
				element.ButtonElement("edit-btn", "button", "Edit", "", html.Data("on-click", datastar.GetSSE("/edit-form/%s", post.ID))),
				element.ButtonElement("view-post", "button", "View More", "", html.Data("on-click", datastar.GetSSE("/post/%s", post.ID))),
			),
		),
	)
}

func AddForm() gomponents.Node {
	return html.Div(
		html.ID(FormId),
		html.Form(
			html.ID("add-form"),
			html.Data("on-submit", "@post('/post', {contentType: 'form', openWhenHidden: true,}); @setAll('input.', '')"),
			html.FieldSet(
				html.Legend(gomponents.Text("Add Post")),
				element.InputElement("Title", "title", "text", "input.title"),
				element.InputElement("Author", "author", "text", "input.author"),
				element.Textarea("Content", "content", "input.content", html.Rows("10")),
				element.ButtonElement("", "submit", "Submit", ""),
			),
		),
	)
}

func EditForm(post *Post) gomponents.Node {
	return html.Div(
		html.ID(FormId),
		html.Data("show", "$view ==='form'"),
		html.Form(
			html.ID(fmt.Sprintf("edit-form-%s", post.ID)),
			html.Data("on-submit", fmt.Sprintf("@put('/post/%s', {contentType: 'form'})", post.ID)),
			html.FieldSet(
				html.Legend(gomponents.Text("Edit Post")),
				element.InputElement("Title", "title", "text", "input.title"),
				element.InputElement("Author", "author", "text", "input.author"),
				element.Textarea("Content", "content", "input.content", html.Rows("10")),
				element.ButtonElement("", "submit", "Submit", ""),
			),
		),
	)
}

func NavbarBackButton() gomponents.Node {
	return components.NavbarComponent(
		NavbarButtonId,
		element.ButtonElement(
			"",
			"button",
			"Back",
			"",
			html.Data("on-click", datastar.GetSSE("/posts")),
		),
	)
}

func NavbarAddPostButton() gomponents.Node {
	return components.NavbarComponent(
		NavbarButtonId,
		element.ButtonElement(
			"",
			"button",
			"Add Post",
			"",
			html.Data("on-click", datastar.GetSSE("/f/add-post")),
		),
	)
}
