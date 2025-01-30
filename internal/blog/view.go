package blog

import (
	"fmt"

	"github.com/reyhardy/go-blog/template/element"
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func PostList(postList Posts) gomponents.Node {
	return html.Div(
		html.ID("postlist"),
		gomponents.Map(postList, func(post *Post) gomponents.Node {
			return PostCard(post)
		},
		))
}

func PostCard(post *Post) gomponents.Node {
	return html.Article(
		html.ID(fmt.Sprintf("post-%s", post.ID)),
		html.Header(html.H1(gomponents.Text(post.Title))),
		html.P(gomponents.Text(post.Content)),
		html.P(html.Cite(gomponents.Text(fmt.Sprintf("- %s", post.Author)))),
		html.Footer(
			element.ButtonElement("button", "Delete", gomponents.Attr("data-on-click", fmt.Sprintf("@delete('/deletepost/%s')", post.ID))),
		),
	)
}

func RenderPostList(postList Posts) gomponents.NodeFunc {
	return PostList(postList).Render
}

func RenderPostCard(post *Post) gomponents.NodeFunc {
	return PostCard(post).Render
}
