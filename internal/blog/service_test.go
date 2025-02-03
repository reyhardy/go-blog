package blog_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/reyhardy/go-blog/internal/blog"
	"github.com/scylladb/gocqlx/v3/qb"
)

func TestAdd(t *testing.T) {
	t.Run("Add Post", func(t *testing.T) {
		postParams := blog.PostParams{
			Title:   "test title",
			Content: "test content",
			Author:  "test author",
		}

		addedPost, err := svc.Add(context.Background(), testKeyspace, &postParams)
		if err != nil {
			t.Errorf("error add: %s", err)
		}

		getPostQuery := session.Query(
			qb.Select(fmt.Sprintf("%v.%v", testKeyspace, blog.TablePost)).
				Columns("title", "content", "author").
				Where(qb.Eq("id")).
				ToCql(),
		).BindStruct(blog.Post{ID: postParams.ID})

		post := new(blog.Post)

		if err = getPostQuery.Scan(
			&post.Title,
			&post.Content,
			&post.Author,
		); err != nil {
			t.Errorf("scan query err: %v", err)
		}

		if addedPost.Title != postParams.Title {
			t.Errorf("added %v got %v", postParams.Title, addedPost.Title)
		}

		if post.Title != postParams.Title {
			t.Errorf("got %v want %v", post.Title, postParams.Title)
		}

		// t.Logf("title postparams: %v, title post: %v, added post title: %v", postParams.Title, post.Title, addedPost.Title)
	})
}

func TestGet(t *testing.T) {
	t.Run("Get Post", func(t *testing.T) {
		_, err := svc.SelectAll(context.Background(), testKeyspace)
		if err != nil {
			t.Errorf("error select all: %s", err)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete Post", func(t *testing.T) {

	})
}
