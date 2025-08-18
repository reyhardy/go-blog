package blog

import (
	"fmt"

	"github.com/reyhardy/go-blog/template/element"
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// post view
func PostList(posts Posts) gomponents.Node {
	return html.Div(
		html.ID("posts"),
		html.Data("on-load", `@get("/api/posts/sse")`),
		gomponents.Map(posts, func(post *Post) gomponents.Node {
			return PostCard(post)
		}),
	)
}

func PostCard(post *Post) gomponents.Node {
	return html.Article(
		html.ID(fmt.Sprintf("post-%s", post.ID)),
		html.Header(
			html.H3(
				html.ID("post_title"),
				gomponents.Textf("Title: %s", post.Title),
			),
		),
		html.P(
			html.ID("post_content"),
			gomponents.Text(post.Content),
		),
		html.Footer(
			html.P(
				html.Cite(
					html.ID("post_author"),
					gomponents.Textf("Posted By: %s", post.Author),
				),
			),
			html.P(
				html.Data("ignore-morph", ""),
				html.ID("post_created_at"),
				gomponents.Textf("Created At: %v", post.CreatedAt.Format("January 2, 2006 @ 3:04:05 PM")),
			),
			gomponents.If(
				post.UpdatedAt != post.CreatedAt,
				html.P(
					html.ID("post_updated_at"),
					gomponents.Textf("Updated At: %v", post.UpdatedAt.Format("January 2, 2006 @ 3:04:05 PM")),
				),
			),
			html.Div(
				html.ID("grid"),
				element.ButtonElement(
					"button",
					"Delete Post",
					html.Data("on-click", fmt.Sprintf(`@delete("/api/post/%s")`, post.ID)),
				),
				ModalEditPost(post),
			),
		),
	)
}

func ModalAddPost() gomponents.Node {
	return html.Div(
		html.Dialog(
			html.Data("ref", "_modalAdd"),
			html.Data("on-click", `evt.target === $_modalAdd && $_modalAdd.close(); @setAll("", {include: /^input\./})`),
			html.Article(
				html.Header(
					html.Button(
						html.Aria("label", "Close"),
						html.Rel("prev"),
						html.Data("on-click",
							`$_modalAdd.close();
							@setAll("", {include: /^input\./})
						`),
					),
				),
				html.Form(
					html.ID("form-add-post"),
					html.FieldSet(
						html.Legend(html.H3(gomponents.Text("Add New Post"))),
						element.InputElement("Post Title", "title", "", "text"),
						element.InputElement("Author Name", "author", "", "text"),
						element.Textarea("Post Content", "content", "", 5),
					),
				),
				html.Footer(
					html.Div(
						html.Button(
							html.Data("on-click",
								`@post("/api/posts", {contentType: "form", selector: "#form-add-post"}); 
								$_modalAdd.close();
								@setAll("", {include: /^input\./})
							`),
							html.Type("button"),
							html.Value("Add Post"),
							gomponents.Text("Add Post"),
						),
					),
				),
			),
		),
		element.ButtonElement(
			"submit",
			"Add Post",
			html.Data("on-click", `$_modalAdd.showModal()`),
		),
	)
}

func ModalEditPost(post *Post) gomponents.Node {
	return html.Div(
		html.Dialog(
			html.Data("ref", fmt.Sprintf("_modalEdit_%s", post.ID)),
			html.Data("on-click", fmt.Sprintf(`evt.target === $_modalEdit_%s && $_modalEdit_%s.close(); @setAll("", {include: /^input\./})`, post.ID, post.ID)),
			html.Article(
				html.Header(
					html.Button(
						html.Aria("label", "Close"),
						html.Rel("prev"),
						html.Data("on-click", fmt.Sprintf(`$_modalEdit_%s.close(); @setAll("", {include: /^input\./})`, post.ID)),
					),
				),
				html.Form(
					html.ID(fmt.Sprintf("form-edit-post-%s", post.ID)),
					html.FieldSet(
						html.Legend(html.H3(gomponents.Text("Edit Post"))),
						element.InputElement("Post Title", "title", post.Title, "text"),
						element.InputElement("Author Name", "author", post.Author, "text"),
						element.Textarea("Post Content", "content", post.Content, 5),
					),
				),
				html.Footer(
					html.Div(
						html.Button(
							html.Data("on-click",
								fmt.Sprintf(`@put("/api/post/%s", {contentType: "form", selector: "#form-edit-post-%s"}); 
								$_modalEdit_%s.close();
								@setAll("", {include: /^input\./})
							`, post.ID, post.ID, post.ID)),
							html.Type("button"),
							html.Value("Edit Post"),
							gomponents.Text("Edit Post"),
						),
					),
				),
			),
		),
		element.ButtonElement(
			"button",
			"Edit Post",
			html.Data("on-click", fmt.Sprintf("$_modalEdit_%s.showModal()", post.ID)),
		),
	)
}
